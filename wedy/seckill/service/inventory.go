package service

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/repository"
	"context"
	"errors"
)

type InventoryService interface {
	VerifyInventory(ctx context.Context, tccId string, order domain.Order) error
	ReduceInventory(ctx context.Context, tccId string) error
	Withdraw(ctx context.Context, tccId string) error
	WithHold(ctx context.Context, activityID int64, productId int64, quantity int64) error
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
