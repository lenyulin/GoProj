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

const BizId = "order"

type order struct {
	log        logger.LoggerV1
	orderRepo  repository.OrderRepository
	tccManager TccManagerService
}

func NewOrderService() OrderService {
	return &order{}
}

func (o *order) Create(ctx context.Context, tccId string, order domain.Order) error {
	orderTx := domain.OrderTX{
		OrderId:     order.OrderId,
		UserId:      order.UserId,
		Price:       order.Price,
		Quantity:    order.Quantity,
		PromoteCode: order.PromoCode,
	}
	err := o.tccManager.AddTcc(ctx, orderTx, BizId+tccId)
	if err != nil {
		return err
	}
	_, err = o.orderRepo.Create(ctx, order)
	if err != nil {
		err = o.tccManager.Failed(ctx, BizId+tccId)
		return err
	}
	err = o.tccManager.Succeed(ctx, BizId+tccId)
	return err
}

func (o *order) Cancel(ctx context.Context) error {
	_, err := o.orderRepo.Cancel(ctx)
	return err
}

func (o *order) Status(ctx context.Context, order domain.Order) (string, error) {
	//TODO implement me
	panic("implement me")
}
