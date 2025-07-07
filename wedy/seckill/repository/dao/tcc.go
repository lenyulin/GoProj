package dao

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/goccy/go-json"
	"strconv"
)

const (
	TccTransactionWatchTopic = "seckill_tcc_transaction_watch"
	SubmitCancelRequest      = "SubmitCancelRequest"
	AddTransaction           = "AddTransaction"
	TransactionComplete      = "TransactionComplete"
	TransactionFailed        = "TransactionFailed"
)
const delayDuration = 3000

type Tcc interface {
	Cancel(ctx context.Context, order OrderTXDAO, id string) error
	Add(ctx context.Context, order OrderTXDAO, id string) error
	Complete(ctx context.Context, order OrderTXDAO, id string) error
	Failed(ctx context.Context, order OrderTXDAO, id string) error
}

type tcc struct {
	producer sarama.SyncProducer
}

func (t *tcc) Failed(ctx context.Context, order OrderTXDAO, id string) error {
	return t.sendMessage(order, id, TransactionFailed)
}

func (t *tcc) Cancel(ctx context.Context, order OrderTXDAO, id string) error {
	return t.sendMessage(order, id, SubmitCancelRequest)
}

func (t *tcc) Add(ctx context.Context, order OrderTXDAO, id string) error {
	return t.sendMessage(order, id, AddTransaction)
}

func (t *tcc) Complete(ctx context.Context, order OrderTXDAO, id string) error {
	return t.sendMessage(order, id, TransactionComplete)
}

func (t *tcc) sendMessage(order OrderTXDAO, id string, txStaus string) error {
	o, err := json.Marshal(order)
	if err != nil {
		return err
	}
	_, _, err = t.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TccTransactionWatchTopic,
		Value: sarama.StringEncoder(o),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("tx_id"),
				Value: []byte(id),
			},
			{
				Key:   []byte("tx_staus"),
				Value: []byte(txStaus),
			},
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
