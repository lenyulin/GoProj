package event

import (
	"GoProj/wedy/comment/repository"
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/pkg/saramax"
	"GoProj/wedy/pkg/tccx"
	"context"
	"github.com/IBM/sarama"
)

const (
	TccTransactionWatchTopic = "seckill_tcc_transaction_watch"
	SubmitCancelRequest      = "SubmitCancelRequest"
	AddTransaction           = "AddTransaction"
	TransactionComplete      = "TransactionComplete"
	TransactionFailed        = "TransactionFailed"
)

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
	cg, err := sarama.NewConsumerGroupFromClient("Seckill_Order", i.client)
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
	//for _, header := range msg.Headers {
	//	if string(header.Key) == "ttl" {
	//		ttl, _ = strconv.ParseInt(string(header.Value), 10, 64)
	//	}
	//}
	//txID, txStatus, ttl := string(msg.Value[0]), string(msg.Value[1]), int64(msg.Value[2])
	//switch txStatus {
	//case TransactionFailed:
	//case SubmitCancelRequest:
	//case AddTransaction:
	//case TransactionComplete:
	//}
	//TODO Implement me
	panic("implement me")
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
