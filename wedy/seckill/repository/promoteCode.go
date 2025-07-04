package repository

import "context"

type PromoteCode interface {
	VerifyCode(ctx context.Context, activityId int64, productId int64, code []string) error
	WithHold(ctx context.Context, activityId int64, productId int64, quality int64) error
	Withdraw(ctx context.Context, productId int64, amount int64) error
}
