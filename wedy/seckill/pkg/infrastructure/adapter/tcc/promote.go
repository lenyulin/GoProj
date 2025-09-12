package tcc

import (
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/service"
	"context"
	"errors"
	"fmt"
)

type CouponAdapter struct {
	couponService service.CouponService // 依赖优惠券服务
}

func NewCouponAdapter(couponService service.CouponService) *CouponAdapter {
	return &CouponAdapter{
		couponService: couponService,
	}
}

func (c *CouponAdapter) Try(ctx context.Context, tccID string, data interface{}) error {
	order, ok := data.(domain.Order)
	if !ok {
		return errors.New("invalid biz data type, expect Order")
	}
	if err := c.couponService.VerifyCode(ctx, tccID, order); err != nil {
		return fmt.Errorf("coupon verify failed: %v", err)
	}
	return c.couponService.LockCoupon(
		ctx,
		order.PromoCode,
		order.OrderId,
	)
}

func (c *CouponAdapter) Confirm(ctx context.Context, tccID string) error {
	return c.couponService.UseCoupon(ctx, tccID)
}

func (c *CouponAdapter) Cancel(ctx context.Context, tccID string) error {
	// 实际项目中需要通过tccID查询关联的优惠券信息
	return c.couponService.UnlockCoupon(ctx, tccID)
}
