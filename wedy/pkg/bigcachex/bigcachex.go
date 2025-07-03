package bigcachex

import "context"

type BigCachex interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
	Delete(ctx context.Context, key string) error
	Serve() error
	Shutdown() error
}
