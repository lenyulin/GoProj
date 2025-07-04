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

func (s *seckill) tccIds(activityId, productId, userId int64) string {
	return fmt.Sprintf("%d:%d:%d:%d:%d", s.biz, s.bizId, activityId, productId, userId)
}

func (s *seckill) Processing(ctx context.Context, order domain.Order) (string, error) {
	var wg sync.WaitGroup
	ch := make(chan error, 2)
	defer close(ch)
	var err error
	//tccId := s.tccIds(order.ActivityId, order.ProductId, order.UserId)
	//Try
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
			_ = s.tccManage.CheckTcc(ctx, s.tccIds(order.OrderId, order.ProductId, order.UserId))
		}
	}
	orderId, err := s.orderSvc.Commit(ctx, order)
	if err != nil {
		return "", err
	}
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
