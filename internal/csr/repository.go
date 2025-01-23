package csr

import (
	"context"
	"time"

	"github.com/android-sms-gateway/ca/pkg/client"
	"github.com/redis/go-redis/v9"
)

type repository struct {
	redis *redis.Client

	prefix string
	ttl    time.Duration
}

func (r *repository) Create(ctx context.Context, requestId string, csr CSR) error {
	key := r.prefix + ":csr:" + requestId
	keyStatus := r.prefix + ":csr:" + requestId + ":status"

	res := r.redis.SetNX(ctx, key, csr.Content, r.ttl)
	if err := res.Err(); err != nil {
		return err
	}

	if !res.Val() {
		return ErrCSRAlreadyExists
	}

	return r.redis.Set(ctx, keyStatus, client.CSRStatusPending, r.ttl).Err()
}

func (r *repository) GetStatus(ctx context.Context, requestId string) (CSRStatus, error) {
	keyStatus := r.prefix + ":csr:" + requestId + ":status"
	res := r.redis.Get(ctx, keyStatus)
	if err := res.Err(); err != nil {
		return CSRStatus{}, err
	}

	return CSRStatus{
		Status: client.CSRStatus(res.Val()),
	}, nil
}

func newRepository(redis *redis.Client, prefix string, ttl time.Duration) *repository {
	if redis == nil {
		panic("redis is required")
	}

	return &repository{
		redis: redis,

		prefix: prefix,
		ttl:    ttl,
	}
}
