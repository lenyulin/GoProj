package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type PromoteCache interface {
	VerifyCode(ctx context.Context, couponIds []string, userId int64) error
	LockCoupon(ctx context.Context, couponIds []string, tcc string) error
	UseCoupon(ctx context.Context, tccId string) error
	UnlockCoupon(ctx context.Context, tccId string) error
}
type promoteRedisCache struct {
	cache *redis.Cmdable
}
