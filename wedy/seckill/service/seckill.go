package service

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/seckill/domain"
	"context"
	"fmt"
	"strconv"
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
	tccManage    TccManagerService
	orderSvc     OrderService
	promoteSvc   PromoteService
	inventorySvc InventoryService
}

func (s *seckill) tccIds(activityId, productId, userId int64) string {
	return fmt.Sprintf(":%d:%d:%d:%d:%d", s.biz, s.bizId, activityId, productId, userId)
}

func (s *seckill) Processing(ctx context.Context, order domain.Order) (string, error) {
	var wg sync.WaitGroup
	ch := make(chan error, 3)
	defer close(ch)
	tccId := s.tccIds(order.OrderId, order.ProductId, order.UserId)
	//Try
	_ = s.tccManage.AddTcc(ctx, domain.OrderTX{
		OrderId:     order.OrderId,
		UserId:      order.UserId,
		Price:       order.Price,
		PromoteCode: order.PromoCode,
	}, "seckill"+tccId)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- s.promoteSvc.VerifyCode(ctx, tccId, order)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- s.inventorySvc.ReduceInventory(ctx, tccId, order)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- s.orderSvc.Create(ctx, tccId, order)
	}()
	wg.Wait()
	//检查事务是否失败
	for er := range ch {
		if er != nil {
			_ = s.tccManage.Failed(ctx, "seckill"+tccId)
			return "", er
		}
	}
	_ = s.tccManage.Succeed(ctx, "seckill"+tccId)
	return strconv.FormatInt(order.OrderId, 10), nil
}

func (s *seckill) Cancel(ctx context.Context) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *seckill) Status(ctx context.Context, order domain.Order) (string, error) {
	//TODO implement me
	panic("implement me")
}
