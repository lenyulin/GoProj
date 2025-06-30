package video

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

const TopicWatchEvent = "video_watch"

type VideoProducer interface {
	VideoProduceWatchEvent(evt WatchEvent) error
}
type WatchEvent struct {
	//video ID
	Vid int64
	//Reader ID
	Uid int64
}
type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) VideoProducer {
	return &SaramaSyncProducer{
		producer: producer,
	}
}

func (s *SaramaSyncProducer) VideoProduceWatchEvent(evt WatchEvent) error {
	val, err := json.Marshal(&evt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicWatchEvent,
		Value: sarama.StringEncoder(val),
	})
	return err
}
