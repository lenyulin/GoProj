package bigcachex

import "context"

type HybridCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
	Delete(ctx context.Context, key string) error
	Serve() error
	Shutdown() error
}
