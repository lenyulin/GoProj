package sarama

import (
	"fmt"
	"github.com/IBM/sarama"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

var addr = []string{"14.103.175.18:9094"}

func TestKafkaSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addr, cfg)
	assert2.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("hello world"),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
		},
		Metadata: "meta data",
	})
}
func TestKafkaAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(addr, cfg)
	assert2.NoError(t, err)
	msgs := producer.Input()
	msgs <- &sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("hello world"),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
		},
		Metadata: "meta data",
	}
	select {
	case msg := <-producer.Successes():
		fmt.Println("send ok.", string(msg.Value.(sarama.StringEncoder)))
	case err := <-producer.Errors():
		fmt.Println("send faild.", err.Error())
	}
}
