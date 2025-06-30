package sarama

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addr, "demo", cfg)
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var consumerhdl ConsumerHdl
	err = consumer.Consume(ctx, []string{"test_topic"}, consumerhdl)
	assert.NoError(t, err)
}

type ConsumerHdl struct {
}

func (c ConsumerHdl) Setup(session sarama.ConsumerGroupSession) error {
	var offset int64 = 0
	partitions := session.Claims()["test_topic"]
	for _, partition := range partitions {
		session.ResetOffset("test_topic", partition, offset, "")
	}
	return nil
}

func (c ConsumerHdl) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c ConsumerHdl) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchsize = 10
	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchsize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		var done = false
		var eg errgroup.Group
		for i := 0; i < batchsize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				batch = append(batch, msg)
				eg.Go(func() error {
					//并发处理错误
					fmt.Printf("Consumer Message:\n%s\n", string(msg.Value))
					return nil
				})
			}
		}
		cancel()
		err := eg.Wait()
		if err != nil {
			fmt.Printf("ERROR")
			continue
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
