package service

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/events"
	"GoProj/wedy/seckill/pkg/tcc"
	"context"
	"fmt"
	"strconv"
	"time"
)

type Seckill interface {
	Processing(ctx context.Context, order domain.Order) (string, error)
	Cancel(ctx context.Context) (string, error)
	Status(ctx context.Context, order domain.Order) (string, error)
}

type seckill struct {
	biz      int64
	bizId    int64
	log      logger.LoggerV1
	adaptors []tcc.TccAction
	//stockAdapter  *tcc3.InventoryAdapter
	//couponAdapter *tcc3.CouponAdapter
	//orderAdapter  *tcc3.OrderAdapter
	tccManager tcc.TCCManagerV1
	producer   events.TCCMegProducer
}

func (s *seckill) tccIds(orderId int64) string {
	return fmt.Sprintf(":%d:%d:%d", s.biz, s.bizId, orderId)
}
func NewSeckillSvc(biz int64, bizId int64, log logger.LoggerV1, adaptors []tcc.TccAction) Seckill {
	return &seckill{
		biz:      biz,
		bizId:    bizId,
		log:      log,
		adaptors: adaptors,
	}
}
func (m *seckill) add(gtid string, Data interface{}) error {
	id, err := strconv.ParseInt(gtid, 10, 64)
	if err != nil {
		return err
	}
	err = m.producer.TCCMangerProduceAddTCCEvent(events.AddTCCEvent{
		TCCIdx:    id,
		TimeStamp: time.Now().UnixMilli(),
		DATA:      Data,
	})
	return err
}
func (s *seckill) Processing(ctx context.Context, order domain.Order) (string, error) {
	gtid := s.tccIds(order.OrderId)
	err := s.add(gtid, order)
	if err != nil {
		return "", err
	}
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
