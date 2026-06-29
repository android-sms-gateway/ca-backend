package csr

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"time"

	"github.com/android-sms-gateway/client-go/ca"
	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
	"github.com/jaevor/go-nanoid"
	"go.uber.org/zap"
)

const serialNumberLength = 120

type Service struct {
	csrs *repository

	caCert *x509.Certificate
	caKey  any

	queue *queue.Queue
	newid func() string
	log   *zap.Logger
}

func NewService(csrs *repository, caCert *x509.Certificate, caKey any, log *zap.Logger) *Service {
	if csrs == nil {
		panic("csrs is required")
	}

	if caCert == nil {
		panic("caCert is required")
	}

	if caKey == nil {
		panic("caKey is required")
	}

	if log == nil {
		panic("log is required")
	}

	newid, _ := nanoid.Canonic()

	s := &Service{
		csrs: csrs,

		caCert: caCert,
		caKey:  caKey,

		queue: nil,
		newid: newid,
		log:   log,
	}

	s.queue = queue.NewPool(
		int64(runtime.GOMAXPROCS(0)),
		queue.WithFn(s.process),
		queue.WithLogger(log.Sugar()),
	)

	return s
}

func (s *Service) Create(ctx context.Context, csr CSR) (Status, error) {
	req, err := s.parseCsr(csr.content)
	if err != nil {
		s.log.Error("failed to parse csr", zap.Error(err))
		return Status{}, fmt.Errorf("%w: %w", ErrCSRInvalid, err)
	}

	if len(req.IPAddresses) != 1 {
		s.log.Error("invalid csr", zap.Any("csr", req))
		return Status{}, fmt.Errorf("%w: should have exactly one IP address", ErrCSRInvalid)
	}

	if req.Subject.CommonName != req.IPAddresses[0].String() {
		s.log.Error("invalid csr", zap.Any("csr", req))
		return Status{}, fmt.Errorf("%w: common name and IP address should be the same", ErrCSRInvalid)
	}

	if !req.IPAddresses[0].IsPrivate() {
		s.log.Error("invalid csr", zap.Any("csr", req))
		return Status{}, fmt.Errorf("%w: IP address should be private", ErrCSRInvalid)
	}

	id := s.newid()
	if insErr := s.csrs.Insert(ctx, id, csr); insErr != nil {
		s.log.Error("failed to create csr", zap.Error(insErr))
		return Status{}, insErr
	}

	if queueErr := s.queue.Queue(csrID(id)); queueErr != nil {
		s.log.Error("failed to queue csr", zap.Error(queueErr))
	}

	return NewCSRStatus(id, csr.csrType, csr.content, csr.metadata, ca.CSRStatusPending, "", ""), nil
}

func (s *Service) Get(ctx context.Context, id string) (Status, error) {
	return s.csrs.Get(ctx, id)
}

func (s *Service) Stop(_ context.Context) error {
	s.queue.Release()

	return nil
}

func (s *Service) process(ctx context.Context, m core.TaskMessage) error {
	id := string(m.Payload())

	res, err := s.csrs.Get(ctx, id)
	if err != nil {
		return err
	}

	if res.Status() != ca.CSRStatusPending {
		return nil
	}

	csr, err := s.parseCsr(res.Content())
	if err != nil {
		return err
	}

	prefix, ok := csrTypeToPrefix[res.csrType]
	if !ok {
		return fmt.Errorf("%w: unknown csr type: %s", ErrCSRInvalid, res.csrType)
	}

	serialNumber, err := s.newSerialNumber(prefix)
	if err != nil {
		return err
	}

	// Create a signed certificate
	template := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               csr.Subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		Extensions:            csr.Extensions,
		ExtraExtensions:       csr.ExtraExtensions,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		EmailAddresses:        csr.EmailAddresses,
		IPAddresses:           csr.IPAddresses,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, s.caCert, csr.PublicKey, s.caKey)
	if err != nil {
		return fmt.Errorf("failed to sign certificate: %w", err)
	}

	// Encode the signed certificate to PEM format
	var certPEM strings.Builder
	if encErr := pem.Encode(&certPEM, &pem.Block{Type: "CERTIFICATE", Headers: nil, Bytes: certBytes}); encErr != nil {
		return fmt.Errorf("failed to encode certificate: %w", encErr)
	}

	s.log.Info(
		"signed certificate",
		zap.String("id", id),
		zap.String("csr", res.Certificate()),
		zap.String("cert", certPEM.String()),
	)

	return s.csrs.SetCertificate(ctx, id, certPEM.String())
}

func (s *Service) parseCsr(content string) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode([]byte(content))
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("%w: can't decode PEM block or invalid block type", ErrCSRInvalid)
	}

	req, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse csr: %w", err)
	}

	return req, nil
}

func (s *Service) newSerialNumber(prefix SerialNumberPrefix) (*big.Int, error) {
	serialNumberLimit := new(big.Int).
		Lsh(big.NewInt(1), serialNumberLength)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}
	return serialNumber.Or(
		serialNumber,
		new(big.Int).Lsh(big.NewInt(int64(prefix)), serialNumberLength),
	), nil
}
