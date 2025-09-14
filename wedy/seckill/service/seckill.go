package service

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/seckill/domain"
	event "GoProj/wedy/seckill/events"
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
	producer event.TCCMegProducer
}

func (s *seckill) tccIds(orderId int64) string {
	return fmt.Sprintf(":%d:%d:%d", s.biz, s.bizId, orderId)
}

func NewSeckillSvc(biz int64, bizId int64, log logger.LoggerV1) Seckill {
	return &seckill{
		biz:   biz,
		bizId: bizId,
		log:   log,
	}
}
func (m *seckill) add(ctx context.Context, gtid string, Data interface{}) error {
	err := m.producer.TCCMangerProduceAddTCCEvent(event.AddTCCEvent{
		TCCIdx:    gtid,
		Status:    "TRYING",
		TimeStamp: time.Now().UnixMilli(),
		DATA:      Data,
	})
	return err
}
func (s *seckill) Processing(ctx context.Context, order domain.Order) (string, error) {
	gtid := s.tccIds(order.OrderId)
	err := s.add(ctx, gtid, order)
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
