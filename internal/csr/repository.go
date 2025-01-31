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

func (r *repository) Insert(ctx context.Context, requestId string, csr CSR) error {
	res := r.redis.HSetNX(ctx, keyStatus, requestId, string(ca.CSRStatusPending))
	if err := res.Err(); err != nil {
		return fmt.Errorf("failed to create csr: %w", err)
	}

	if !res.Val() {
		return ErrCSRAlreadyExists
	}

	key := "csr:" + requestId
	validUntil := time.Now().Add(r.ttl)
	_, err := r.redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, csr.toMap())
		pipe.ExpireAt(ctx, key, validUntil)
		pipe.HExpireAt(ctx, keyStatus, validUntil, requestId)

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create csr: %w", errors.Join(err, r.redis.HDel(ctx, keyStatus, requestId).Err()))
	}

	return nil
}

func (r *repository) Get(ctx context.Context, requestId string) (CSRStatus, error) {
	status, err := r.redis.HGet(ctx, keyStatus, requestId).Result()
	if errors.Is(err, redis.Nil) {
		return CSRStatus{}, ErrCSRNotFound
	}
	if err != nil {
		return CSRStatus{}, fmt.Errorf("failed to get csr: %w", err)
	}

	key := "csr:" + requestId
	res, err := r.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return CSRStatus{}, fmt.Errorf("failed to get csr: %w", err)
	}

	if len(res) == 0 {
		return CSRStatus{}, ErrCSRNotFound
	}

	metadata := map[string]string{}

	if err := json.Unmarshal([]byte(res["metadata"]), &metadata); err != nil {
		return CSRStatus{}, fmt.Errorf("failed to get csr: %w", err)
	}

	return NewCSRStatus(requestId, res["content"], metadata, ca.CSRStatus(status), res["certificate"], res["reason"]), nil
}

func (r *repository) SetCertificate(ctx context.Context, requestId string, certificate string) error {
	key := "csr:" + requestId

	_, err := r.redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, "certificate", certificate)
		pipe.HSet(ctx, keyStatus, requestId, string(ca.CSRStatusApproved))

		return nil
	})

	return err
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
