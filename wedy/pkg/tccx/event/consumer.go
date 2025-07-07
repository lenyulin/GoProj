package event

import (
	"GoProj/wedy/comment/repository"
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/pkg/saramax"
	"GoProj/wedy/pkg/tccx"
	"context"
	"github.com/IBM/sarama"
	"strconv"
)

const TccTransactionWatchTopic = "seckill_tcc_transaction_watch"

type CommentWatchEventConsumer struct {
	client  sarama.Client
	manager tccx.TccManager
	log     logger.LoggerV1
}

func NewInteractiveWatchEventConsumer(repo repository.CommentRepository, client sarama.Client, log logger.LoggerV1) *CommentWatchEventConsumer {
	return &CommentWatchEventConsumer{
		client: client,
		log:    log,
	}
}
func (i *CommentWatchEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("comment", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TccTransactionWatchTopic},
			saramax.NewHandler[tccx.TccTransactionEvent](i.log, i.Consume))
		if er != nil {
			//转发消息到死信队列
			i.log.Error("Out off Consumer Error @wedy/comment/events/consumer.go line:43", logger.Error(er))
		}
	}()
	return nil
}

func (i *CommentWatchEventConsumer) Consume(msg *sarama.ConsumerMessage, event tccx.TccTransactionEvent) error {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//txID, txStatus, ttl := i.parseHeaders(msg.Headers)
	//return i.manager.Process(ctx, ,,,event.Order)
	panic("TODO")
}
func (i *CommentWatchEventConsumer) parseHeaders(headers []sarama.RecordHeader) (txID, txStatus string, ttl int) {
	for _, header := range headers {
		key := string(header.Key)
		value := string(header.Value)

		switch key {
		case "tx_id":
			txID = value
		case "tx_staus": // 注意：生产者代码中这里有拼写错误，应该是"tx_status"
			txStatus = value
		case "ttl":
			ttl, _ = strconv.Atoi(value) // 忽略转换错误，默认为0
		}
	}

	return
}

//func (i *CommentWatchEventConsumer) StartBatch() error {
//	cg, err := sarama.NewConsumerGroupFromClient("comment", i.client)
//	if err != nil {
//		return err
//	}
//	go func() {
//		er := cg.Consume(context.Background(),
//			[]string{TopicCommentSubmitEvent},
//			saramax.NewBatchHandler[CommentEvent](i.log, i.BatchConsume))
//		if er != nil {
//			i.log.Error("Out off Consumer Error @wedy/comment/events/consumer.go line:58", logger.Error(er))
//		}
//	}()
//	return nil
//}

//	func (i *CommentWatchEventConsumer) BatchConsume(msg []*sarama.ConsumerMessage, events []CommentEvent) error {
//		bizs := make([]string, 0, len(events))
//		bizIds := make([]int64, 0, len(events))
//		for _, evt := range events {
//			bizs = append(bizs, "comment")
//			bizIds = append(bizIds, evt.Id)
//		}
//		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//		defer cancel()
//		return i.repo.Submit(ctx, bizs, bizIds)
//	}
