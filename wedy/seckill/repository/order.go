package repository

import (
	"GoProj/wedy/seckill/domain"
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, order domain.Order) (string, error)
	Cancel(ctx context.Context) (string, error)
	Status(ctx context.Context, order domain.Order) (string, error)
	Withdraw(ctx context.Context, order domain.Order) (string, error)
}
type orderRepository struct {
}

func (o *orderRepository) Create(ctx context.Context, order domain.Order) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (o *orderRepository) Cancel(ctx context.Context) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (o *orderRepository) Status(ctx context.Context, order domain.Order) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (o *orderRepository) Withdraw(ctx context.Context, order domain.Order) (string, error) {
	//TODO implement me
	panic("implement me")
}
