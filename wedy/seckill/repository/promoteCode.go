package repository

import (
	"GoProj/wedy/seckill/repository/cache"
	"context"
	"errors"
)

type PromoteCode interface {
	VerifyCode(ctx context.Context, couponIds []string, userId int64) error
	LockCoupon(ctx context.Context, couponIds []string, tccId string) error
	UseCoupon(ctx context.Context, tccId string) error
	UnlockCoupon(ctx context.Context, tccId string) error
}

var (
	ErrPromoteCodeNotFound = errors.New("promote code not found")
	ErrPromoteCodeNotMatch = errors.New("promote code not match")
	ErrPromoteCodeExpired  = errors.New("promote code expired")
)

type promoteRepo struct {
	cache cache.PromoteCache
}

func (p *promoteRepo) VerifyCode(ctx context.Context, couponIds []string, userId int64) error {
	return p.cache.VerifyCode(ctx, couponIds, userId)
}

func (p *promoteRepo) LockCoupon(ctx context.Context, couponIds []string, tccId string) error {
	return p.cache.LockCoupon(ctx, couponIds, tccId)
}

func (p *promoteRepo) UseCoupon(ctx context.Context, tccId string) error {
	return p.cache.UseCoupon(ctx, tccId)
}

func (p *promoteRepo) UnlockCoupon(ctx context.Context, tccId string) error {
	return p.cache.UnlockCoupon(ctx, tccId)
}
