package dao

import (
	"context"
	"github.com/IBM/sarama"
)

const SubmitCancelRequestTopic = "SubmitCancelRequestTopic"

type Tcc interface {
	SubmitCancelRequest(ctx context.Context, id string) error
}

type tcc struct {
	producer sarama.SyncProducer
}

func (t *tcc) SubmitCancelRequest(ctx context.Context, id string) error {
	_, _, err := t.producer.SendMessage(&sarama.ProducerMessage{
		Topic: SubmitCancelRequestTopic,
		Value: sarama.StringEncoder("hello world"),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
		},
		Metadata: "meta data",
	})
	return err
}

func NewTccSaramaDAO(producer sarama.SyncProducer) Tcc {
	return &tcc{
		producer: producer,
	}
}
