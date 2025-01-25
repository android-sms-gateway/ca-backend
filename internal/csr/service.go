package csr

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/android-sms-gateway/ca/pkg/client"
	"github.com/jaevor/go-nanoid"
	"go.uber.org/zap"
)

type Service struct {
	csrs *repository

	newid func() string
	log   *zap.Logger
}

func (s *Service) Create(ctx context.Context, csr CSR) (CSRStatus, error) {
	block, _ := pem.Decode([]byte(csr.Content()))
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		s.log.Error("invalid csr", zap.String("csr", csr.Content()))
		return CSRStatus{}, fmt.Errorf("%w: should be a certificate request", ErrCSRInvalid)
	}

	req, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		s.log.Error("failed to parse csr", zap.Error(err))
		return CSRStatus{}, fmt.Errorf("%w: %s", ErrCSRInvalid, err)
	}

	if len(req.IPAddresses) != 1 {
		s.log.Error("invalid csr", zap.Any("csr", req))
		return CSRStatus{}, fmt.Errorf("%w: should have exactly one IP address", ErrCSRInvalid)
	}

	if req.Subject.CommonName != req.IPAddresses[0].String() {
		s.log.Error("invalid csr", zap.Any("csr", req))
		return CSRStatus{}, fmt.Errorf("%w: common name and IP address should be the same", ErrCSRInvalid)
	}

	if !req.IPAddresses[0].IsPrivate() {
		s.log.Error("invalid csr", zap.Any("csr", req))
		return CSRStatus{}, fmt.Errorf("%w: IP address should be private", ErrCSRInvalid)
	}

	id := s.newid()
	if err := s.csrs.Insert(ctx, id, csr); err != nil {
		s.log.Error("failed to create csr", zap.Error(err))
		return CSRStatus{}, err
	}

	return NewCSRStatus(id, client.CSRStatusPending, "", ""), nil
}

func (s *Service) Get(ctx context.Context, id string) (CSRStatus, error) {
	return s.csrs.Get(ctx, id)
}

func NewService(csrs *repository, log *zap.Logger) *Service {
	if csrs == nil {
		panic("csrs is required")
	}

	if log == nil {
		panic("log is required")
	}

	newid, _ := nanoid.Canonic()

	return &Service{
		csrs: csrs,

		newid: newid,
		log:   log,
	}
}
