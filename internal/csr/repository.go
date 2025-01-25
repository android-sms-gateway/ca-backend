package csr

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/android-sms-gateway/ca/pkg/client"
	"github.com/redis/go-redis/v9"
)

const (
	keyStatus = "csr"
)

type repository struct {
	redis *redis.Client

	ttl time.Duration
}

func (r *repository) Create(ctx context.Context, requestId string, csr CSR) error {
	res := r.redis.HSetNX(ctx, keyStatus, requestId, string(client.CSRStatusPending))
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

// func (r *repository) GetStatus(ctx context.Context, requestId string) (CSRStatus, error) {
// 	keyStatus := r.prefix + ":csr:" + requestId + ":status"
// 	res := r.redis.Get(ctx, keyStatus)
// 	if err := res.Err(); err != nil {
// 		return CSRStatus{}, err
// 	}

// 	return CSRStatus{
// 		Status: client.CSRStatus(res.Val()),
// 	}, nil
// }

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
