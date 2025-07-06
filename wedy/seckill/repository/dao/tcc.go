package dao

import (
	"context"
	"github.com/IBM/sarama"
	"strconv"
)

const (
	SubmitCancelRequestTopic = "SubmitCancelRequest"
	AddTransactionTopic      = "AddTransaction"
	TransactionCompleteTopic = "TransactionComplete"
	TransactionFailedTopic   = "TransactionFailed"
)
const delayDuration = 3000

type Tcc interface {
	Cancel(ctx context.Context, id string) error
	Add(ctx context.Context, id string) error
	Complete(ctx context.Context, id string) error
	Failed(ctx context.Context, id string) error
}

type tcc struct {
	producer sarama.SyncProducer
}

func (t *tcc) Failed(ctx context.Context, id string) error {
	return t.sendMessage(id, TransactionFailedTopic)
}

func (t *tcc) Cancel(ctx context.Context, id string) error {
	return t.sendMessage(id, SubmitCancelRequestTopic)
}

func (t *tcc) Add(ctx context.Context, id string) error {
	return t.sendMessage(id, AddTransactionTopic)
}

func (t *tcc) Complete(ctx context.Context, id string) error {
	return t.sendMessage(id, TransactionCompleteTopic)
}

func (t *tcc) sendMessage(id string, topic string) error {
	_, _, err := t.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(id),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("ttl"),
				Value: []byte(strconv.Itoa(delayDuration)),
			},
		},
	})
	return err
}
func NewTccSaramaDAO(producer sarama.SyncProducer) Tcc {
	return &tcc{
		producer: producer,
	}
}
