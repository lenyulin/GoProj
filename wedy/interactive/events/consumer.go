package events

import (
	"GoProj/wedy/interactive/repository"
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/pkg/saramax"
	"context"
	"github.com/IBM/sarama"
	"time"
)

const TopicWatchEvent = "video_watch"

type WatchEvent struct {
	//video ID
	Vid int64
	//Reader ID
	Uid int64
}
type InteractiveWatchEventConsumer struct {
	repo   repository.InteractiveRepository
	client sarama.Client
	log    logger.LoggerV1
}

func NewInteractiveWatchEventConsumer(repo repository.InteractiveRepository, client sarama.Client, log logger.LoggerV1) *InteractiveWatchEventConsumer {
	return &InteractiveWatchEventConsumer{
		repo:   repo,
		client: client,
		log:    log,
	}
}
func (i *InteractiveWatchEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicWatchEvent},
			saramax.NewHandler[WatchEvent](i.log, i.Consume))
		if er != nil {
			i.log.Error("Out off Consumer Error @wedy/interactive/events/video/consumer.go line:29", logger.Error(er))
		}
	}()
	return nil
}
func (i *InteractiveWatchEventConsumer) StartBatch() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicWatchEvent},
			saramax.NewBatchHandler[WatchEvent](i.log, i.BatchConsume))
		if er != nil {
			i.log.Error("Out off Consumer Error @wedy/interactive/events/video/consumer.go line:29", logger.Error(er))
		}
	}()
	return nil
}
func (i *InteractiveWatchEventConsumer) BatchConsume(msg []*sarama.ConsumerMessage, events []WatchEvent) error {
	bizs := make([]string, 0, len(events))
	bizIds := make([]int64, 0, len(events))
	for _, evt := range events {
		bizs = append(bizs, "video")
		bizIds = append(bizIds, evt.Uid)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.BatchIncrReadCnt(ctx, bizs, bizIds)
}
func (i *InteractiveWatchEventConsumer) Consume(msg *sarama.ConsumerMessage, event WatchEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.IncrReadCnt(ctx, "video", event.Vid)
}
