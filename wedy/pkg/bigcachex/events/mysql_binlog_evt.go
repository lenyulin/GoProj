package events

import (
	"GoProj/wedy/pkg/bigcachex"
	"GoProj/wedy/pkg/bigcachex/proto"
	"GoProj/wedy/pkg/canalx"
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/pkg/saramax"
	"context"
	"github.com/IBM/sarama"
	proto2 "google.golang.org/protobuf/proto"
	"strconv"
	"time"
)

type MysqlBinlogConsumer interface {
	Start() error
	Consume() error
}

type Consumer struct {
	consumer sarama.Client
	cache    bigcachex.HybridCache
	log      logger.LoggerV1
	groupId  string
	targetDB string
}

func NewMysqlBinlogConsumer(client sarama.Client, log logger.LoggerV1, groupId string, targetDB string) *Consumer {
	return &Consumer{
		consumer: client,
		log:      log,
		groupId:  groupId,
		targetDB: targetDB,
	}
}

func (m *Consumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient(m.groupId, m.consumer)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"seckill_current_activity_binlog"},
			saramax.NewHandler[canalx.Message[proto.SeckillActivity]](m.log, m.Consume),
		)
		if err != nil {
			m.log.Error("消费异常", logger.Error(err))
		}
	}()
	return err
}

func (m *Consumer) Consume(msg *sarama.ConsumerMessage, val canalx.Message[proto.SeckillActivity]) error {
	if val.Table != m.targetDB {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for i := range val.Data {
		data := &val.Data[i]
		//已结束删除缓存
		if data.Status == -1 {
			err := m.cache.Delete(ctx, strconv.FormatInt(data.ID, 10))
			if err != nil {
				return err
			}
		}
		d, err := proto2.Marshal(data)
		if err != nil {
			return err
		}
		err = m.cache.Set(ctx, strconv.FormatInt(data.ProductID, 10), d)
		if err != nil {
			return err
		}
	}
	return nil
}
