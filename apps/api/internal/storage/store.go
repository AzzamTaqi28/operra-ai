package storage

import "context"

type Store interface {
	EnsureBucket(ctx context.Context) error
	Put(ctx context.Context, key, contentType string, data []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}
