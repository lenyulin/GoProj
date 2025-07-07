package service

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/repository"
	"context"
	"errors"
)

type InventoryService interface {
	Confirm(ctx context.Context, activityId int64, productId int64) (int64, error)
	ReduceInventory(ctx context.Context, tccId string, order domain.Order) error
	Withdraw(ctx context.Context, productId int64, amount int64) error
}

var (
	ErrItemNotFound          = errors.New("item not found")
	ErrItemNotAvailable      = errors.New("item not available")
	ErrInsufficientInventory = errors.New("insufficient inventory")
)

type inventory struct {
	log        logger.LoggerV1
	repo       repository.Inventory
	tccManager TccManagerService
}

func (i *inventory) Confirm(ctx context.Context, activityId int64, productId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (i *inventory) ReduceInventory(ctx context.Context, tccId string, order domain.Order) error {
	orderTx := domain.OrderTX{
		OrderId:     order.OrderId,
		UserId:      order.UserId,
		Price:       order.Price,
		Quantity:    order.Quantity,
		PromoteCode: order.PromoCode,
	}
	err := i.tccManager.AddTcc(ctx, orderTx, "inventory"+tccId)
	if err != nil {
		return err
	}
	err = i.repo.ReserveStock(ctx, order.ActivityId, order.ProductId, order.Quantity)
	if err != nil {
		err = i.tccManager.Failed(ctx, "inventory"+tccId)
		return err
	}
	err = i.tccManager.Succeed(ctx, "inventory"+tccId)
	return err
}

func (i *inventory) Withdraw(ctx context.Context, productId int64, amount int64) error {
	return i.repo.Withdraw(ctx, productId, amount)
}
