package csr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/android-sms-gateway/client-go/ca"
	"github.com/redis/go-redis/v9"
)

const (
	keyStatus = "csr"
)

type repository struct {
	redis *redis.Client

	ttl time.Duration
}

func (r *repository) Insert(ctx context.Context, requestID string, csr CSR) error {
	res := r.redis.HSetNX(ctx, keyStatus, requestID, string(ca.CSRStatusPending))
	if err := res.Err(); err != nil {
		return fmt.Errorf("failed to create csr: %w", err)
	}

	if !res.Val() {
		return ErrCSRAlreadyExists
	}

	key := "csr:" + requestID
	validUntil := time.Now().Add(r.ttl)
	_, err := r.redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, csr.toMap())
		pipe.ExpireAt(ctx, key, validUntil)
		pipe.HExpireAt(ctx, keyStatus, validUntil, requestID)

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create csr: %w", errors.Join(err, r.redis.HDel(ctx, keyStatus, requestID).Err()))
	}

	return nil
}

func (r *repository) Get(ctx context.Context, requestID string) (Status, error) {
	status, err := r.redis.HGet(ctx, keyStatus, requestID).Result()
	if errors.Is(err, redis.Nil) {
		return Status{}, ErrCSRNotFound
	}
	if err != nil {
		return Status{}, fmt.Errorf("failed to get csr: %w", err)
	}

	key := "csr:" + requestID
	res, err := r.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return Status{}, fmt.Errorf("failed to get csr: %w", err)
	}

	if len(res) == 0 {
		return Status{}, ErrCSRNotFound
	}

	metadata := map[string]string{}

	if jsonErr := json.Unmarshal([]byte(res["metadata"]), &metadata); jsonErr != nil {
		return Status{}, fmt.Errorf("failed to get csr: %w", jsonErr)
	}

	return NewCSRStatus(
		requestID,
		ca.CSRType(res["type"]),
		res["content"],
		metadata,
		ca.CSRStatus(status),
		res["certificate"],
		res["reason"],
	), nil
}

func (r *repository) SetCertificate(ctx context.Context, requestID string, certificate string) error {
	key := "csr:" + requestID

	_, err := r.redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, "certificate", certificate)
		pipe.HSet(ctx, keyStatus, requestID, string(ca.CSRStatusApproved))

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to set certificate: %w", err)
	}

	return nil
}

func newRepository(redis *redis.Client, ttl time.Duration) *repository {
	if redis == nil {
		panic("redis is required")
	}

	if ttl <= 0 {
		panic("ttl must be greater than 0")
	}

	return &repository{
		redis: redis,

		ttl: ttl,
	}
}
