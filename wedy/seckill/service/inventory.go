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
	ReduceInventory(ctx context.Context, order domain.Order) error
	Withdraw(ctx context.Context, productId int64, amount int64) error
}

var (
	ErrItemNotFound          = errors.New("item not found")
	ErrItemNotAvailable      = errors.New("item not available")
	ErrInsufficientInventory = errors.New("insufficient inventory")
)

type inventory struct {
	log  logger.LoggerV1
	repo repository.Inventory
}

func (i *inventory) Confirm(ctx context.Context, activityId int64, productId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (i *inventory) ReduceInventory(ctx context.Context, activityId int64, productId int64, quality int64) error {
	//try
	//var errs []errorx.TccRes
	//err := i.repo.ReserveStock(ctx, productId, quality)
	//append(errs, errorx.TccRes{})
	err := i.repo.ReserveStock(ctx, productId, quality)
	if err != nil {
		errors.Join(err, i.repo.Withdraw(ctx, productId, quality))
		return err
	}
	return nil
}

func (i *inventory) Withdraw(ctx context.Context, productId int64, amount int64) error {
	return i.repo.Withdraw(ctx, productId, amount)
}
