package service

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/seckill/domain"
	tcc3 "GoProj/wedy/seckill/pkg/infrastructure/adapter/tcc"
	"GoProj/wedy/seckill/pkg/tcc"
	"context"
	"fmt"
	"time"
)

type Seckill interface {
	Processing(ctx context.Context, order domain.Order) (string, error)
	Cancel(ctx context.Context) (string, error)
	Status(ctx context.Context, order domain.Order) (string, error)
}

type seckill struct {
	biz           int64
	bizId         int64
	log           logger.LoggerV1
	stockAdapter  *tcc3.InventoryAdapter
	couponAdapter *tcc3.CouponAdapter
	orderAdapter  *tcc3.OrderAdapter
	tccManager    tcc.TCCManager
}

func (s *seckill) tccIds(orderId int64) string {
	return fmt.Sprintf(":%d:%d:%d", s.biz, s.bizId, orderId)
}

func (s *seckill) Processing(ctx context.Context, order domain.Order) (string, error) {
	gtid := s.tccIds(order.OrderId)
	_ = s.tccManager.RegisterAction(gtid, s.stockAdapter)
	_ = s.tccManager.RegisterAction(gtid, s.couponAdapter)
	_ = s.tccManager.RegisterAction(gtid, s.orderAdapter)
	tx := s.tccManager.NewTransaction(gtid, order)
	c, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := s.tccManager.RunTransaction(c, tx); err != nil {
		// 处理事务执行失败
		fmt.Printf("Transaction failed: %v\n", err)
	} else {
		fmt.Println("Transaction completed successfully")
	}
}

func (s *seckill) Cancel(ctx context.Context) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *seckill) Status(ctx context.Context, order domain.Order) (string, error) {
	//TODO implement me
	panic("implement me")
}
