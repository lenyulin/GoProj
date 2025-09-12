package service

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/repository"
	"context"
)

type CouponService interface {
	VerifyCode(ctx context.Context, tccId string, order domain.Order) error
	LockCoupon(ctx context.Context, couponIds []string, orderId int64) error
	UseCoupon(ctx context.Context, tccId string) error
	UnlockCoupon(ctx context.Context, tccId string) error
}

var (
	ErrPromoteCodeNotFound = repository.ErrPromoteCodeNotFound
	ErrPromoteCodeNotMatch = repository.ErrPromoteCodeNotMatch
	ErrPromoteCodeExpired  = repository.ErrPromoteCodeExpired
)

type promote struct {
	log  logger.LoggerV1
	repo repository.PromoteCode
}

func (p *promote) LockCoupon(ctx context.Context, couponIds []string, orderId string) error {
	return p.repo.LockCoupon(ctx, couponIds, orderId)
}

func (p *promote) VerifyCode(ctx context.Context, tccId string, order domain.Order) error {
	return p.repo.VerifyCode(ctx, order.PromoCode, order.UserId)
}

func (p *promote) UseCoupon(ctx context.Context, tccId string) error {
	return p.repo.UseCoupon(ctx, tccId)
}

func (p *promote) UnlockCoupon(ctx context.Context, tccId string) error {
	return p.repo.UnlockCoupon(ctx, tccId)
}
