package repository

import (
	"GoProj/wedy/seckill/domain"
	"context"
)

type PromoteCode interface {
	VerifyCode(ctx context.Context, tccId string, order domain.Order) error
	WithHold(ctx context.Context, activityId string, productId int64, quality int64) error
	Withdraw(ctx context.Context, productId string, amount int64) error
}
