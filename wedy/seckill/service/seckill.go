package service

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/seckill/domain"
	"context"
	"fmt"
	"sync"
)

type Seckill interface {
	Processing(ctx context.Context, order domain.Order) (string, error)
	Cancel(ctx context.Context) (string, error)
	Status(ctx context.Context, order domain.Order) (string, error)
}

type seckill struct {
	biz          int64
	bizId        int64
	log          logger.LoggerV1
	tccManage    TccWithDrawnManage
	orderSvc     OrderService
	promoteSvc   PromoteService
	inventorySvc InventoryService
}

func (s *seckill) tccIds(transaction string, activityId, productId, userId int64) string {
	return fmt.Sprintf("%s:%d:%d:%d:%d:%d", transaction, s.biz, s.bizId, activityId, productId, userId)
}

func (s *seckill) Processing(ctx context.Context, order domain.Order) (string, error) {
	var wg sync.WaitGroup
	ch := make(chan error, 2)
	defer close(ch)
	var err error
	tccId := s.tccIds("Order", order.OrderId, order.ProductId, order.UserId)
	//Try
	_ = s.tccManage.AddTcc(ctx, domain.OrderTX{
		OrderId:     order.OrderId,
		UserId:      order.UserId,
		Mount:       order.Quantity,
		PromoteCode: order.PromoCode,
	}, tccId)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- s.promoteSvc.VerifyCode(ctx, order)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- s.inventorySvc.ReduceInventory(ctx, order)
	}()
	wg.Wait()
	//检查事务是否失败
	for er := range ch {
		if er != nil {
			_ = s.tccManage.Failed(ctx, tccId)
			return "", er
		}
	}
	orderId, err := s.orderSvc.Create(ctx, order)
	if err != nil {
		_ = s.tccManage.Failed(ctx, tccId)
		return "", err
	}
	_ = s.tccManage.Succeed(ctx, tccId)
	return orderId, nil
}

func (s *seckill) Cancel(ctx context.Context) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *seckill) Status(ctx context.Context, order domain.Order) (string, error) {
	//TODO implement me
	panic("implement me")
}
