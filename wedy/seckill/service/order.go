package service

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/repository"
	"context"
)

type OrderService interface {
	Create(ctx context.Context, tccId string, order domain.Order) error
	Cancel(ctx context.Context) error
	Status(ctx context.Context, order domain.Order) (string, error)
}

const serviceClass = "order"

type order struct {
	l         logger.LoggerV1
	orderRepo repository.OrderRepository
}

func NewOrderService(l logger.LoggerV1, orderRepo repository.OrderRepository) OrderService {
	return &order{
		l:         l,
		orderRepo: orderRepo,
	}
}
