package events

import (
	"GoProj/wedy/pkg/canalx"
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/pkg/saramax"
	"context"
	"github.com/IBM/sarama"
	"time"
)

type MysqlBinlogConsumer interface {
	Start() error
	Consume() error
}

type Consumer struct {
	consumer sarama.Client
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
			[]string{"wedy_seckill"},
			saramax.NewHandler[canalx.Message[any]](m.log, m.Consume),
		)
		if err != nil {
			m.log.Error("消费异常", logger.Error(err))
		}
	}()
	return err
}

func (m *Consumer) Consume(msg *sarama.ConsumerMessage, val canalx.Message[any]) error {
	if val.Table != m.targetDB {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, data := range val.Data {
		//TODO
		panic("implement me")
	}
	return nil
}
