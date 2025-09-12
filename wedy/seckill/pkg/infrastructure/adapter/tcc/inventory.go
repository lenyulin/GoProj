package tcc

import (
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/service"
	"context"
	"errors"
	"fmt"
)

type InventoryAdapter struct {
	inventoryService service.InventoryService // 依赖库存服务
}

func NewStockAdapter(inventoryService service.InventoryService) *InventoryAdapter {
	return &InventoryAdapter{
		inventoryService: inventoryService,
	}
}

func (s *InventoryAdapter) Try(ctx context.Context, tccID string, bizData interface{}) error {
	order, ok := bizData.(domain.Order)
	if !ok {
		return errors.New("invalid biz data type, expect Order")
	}
	if err := s.inventoryService.VerifyInventory(ctx, tccID, order); err != nil {
		return fmt.Errorf("stock verify failed: %v", err)
	}
	return s.inventoryService.WithHold(
		ctx,
		order.ActivityId,
		order.ProductId,
		order.Quantity,
	)
}

func (s *InventoryAdapter) Confirm(ctx context.Context, tccID string) error {
	return s.inventoryService.ReduceInventory(ctx, tccID)
}

func (s *InventoryAdapter) Cancel(ctx context.Context, tccID string) error {
	return s.inventoryService.Withdraw(ctx, tccID)
}
