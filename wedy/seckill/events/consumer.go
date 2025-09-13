package events

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/pkg/saramax"
	"GoProj/wedy/seckill/pkg/tcc"
	"context"
	"github.com/IBM/sarama"
	"time"
)

type TCCManagerWatchEventConsumer struct {
	tccManager *tcc.TCCManagerV1
	client     sarama.Client
	log        logger.LoggerV1
}

func NewTCCManagerWatchEventConsumer(tccManager *tcc.TCCManagerV1, client sarama.Client, log logger.LoggerV1) *TCCManagerWatchEventConsumer {
	return &TCCManagerWatchEventConsumer{
		tccManager: tccManager,
		client:     client,
		log:        log,
	}
}

func (i *TCCManagerWatchEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("tcc_manager_group", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicWatchEvent},
			saramax.NewHandler[AddTCCEvent](i.log, i.Consume))
		if er != nil {
			i.log.Error("TCC Consume ERROR", logger.Error(er))
		}
	}()
	return nil
}

func (i *TCCManagerWatchEventConsumer) Consume(msg *sarama.ConsumerMessage, event AddTCCEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return i.tccManager.AddNewTransactionTask(ctx, event.TCCIdx, event.DATA)
}
