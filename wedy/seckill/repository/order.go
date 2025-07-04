package repository

import (
	"GoProj/wedy/seckill/domain"
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, order domain.Order) (string, error)
	Cancel(ctx context.Context) (string, error)
	Status(ctx context.Context, order domain.Order) (string, error)
	Preorder(ctx context.Context, order domain.Order) (string, error)
	Withdraw(ctx context.Context, order domain.Order) (string, error)
}
