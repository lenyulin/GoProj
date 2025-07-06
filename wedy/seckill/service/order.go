package service

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/repository"
	"context"
)

type OrderService interface {
	Create(ctx context.Context, order domain.Order) (string, error)
	Cancel(ctx context.Context) (string, error)
	Status(ctx context.Context, order domain.Order) (string, error)
}

type order struct {
	log       logger.LoggerV1
	orderRepo repository.OrderRepository
}

func NewOrderService() OrderService {
	return &order{}
}

func (o *order) Create(ctx context.Context, order domain.Order) (string, error) {
	panic("implement me")
}

func (o *order) Cancel(ctx context.Context) (string, error) {
	return o.orderRepo.Cancel(ctx)
}

func (o *order) Status(ctx context.Context, order domain.Order) (string, error) {
	//TODO implement me
	panic("implement me")
}
