package repository

import (
	"GoProj/wedy/internal/repository/cache"
	"context"
)

var (
	ErrCodeSendTooMany       = cache.ErrCodeSendTooMany
	ErrCodeVerifySendTooMany = cache.ErrCodeVerifySendTooMany
)

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}
type CachedCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(cache cache.CodeCache) CodeRepository {
	return &CachedCodeRepository{cache}
}
func (c *CachedCodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}
func (c *CachedCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}
