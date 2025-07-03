package service

import (
	"GoProj/wedy/seckill-system/domain"
	"context"
)

type OrderService interface {
	Processing(ctx context.Context) error
	Cancel(ctx context.Context) error
	Status(ctx context.Context, order domain.Order) error
}

type order struct {
}

func NewOrderService() OrderService {
	return &order{}
}
func (o *order) Processing(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (o *order) Cancel(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (o *order) Status(ctx context.Context, order domain.Order) error {
	//TODO implement me
	panic("implement me")
}
