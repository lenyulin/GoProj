package events

import (
	"GoProj/wedy/seckill/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"strconv"
	"time"
)

type TCCMegProducer interface {
	TCCMangerProduceAddTCCEvent(evt AddTCCEvent) error
}

const (
	MaxTCCMangerProduceRetry = 3
)

var (
	ErrProduceSubmitFailure = errors.New("failed to produce summit event")
)

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

var PartitionCount = 3

func (s *SaramaSyncProducer) TCCMangerProduceAddTCCEvent(evt AddTCCEvent) error {
	val, err := json.Marshal(&evt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	evt.TimeStamp = time.Now().UnixMilli()
	partition, offset, err := s.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     TopicWatchEvent,
		Partition: int32(utils.GetPartition(strconv.FormatInt(evt.TCCIdx, 10), PartitionCount)),
		Value:     sarama.StringEncoder(val),
	})
	for err != nil {
		evt.Partition = partition
		evt.Topic = TopicWatchEvent
		evt.Offset = strconv.Itoa(int(offset))
		evt.Retry += 1
		if evt.Retry == MaxTCCMangerProduceRetry {
			//发送到死信队列
			return ErrProduceSubmitFailure
		}
		err = s.TCCMangerProduceAddTCCEvent(evt)
	}
	return nil
}
func (s *SaramaSyncProducer) addTCCConfirmTaskEvent(evt AddTCCEvent) error {
	val, err := json.Marshal(&evt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	evt.TimeStamp = time.Now().UnixMilli()
	partition, offset, err := s.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     TCCConfirmTaskEvent,
		Partition: int32(utils.GetPartition(strconv.FormatInt(evt.TCCIdx, 10), PartitionCount)),
		Value:     sarama.StringEncoder(val),
	})
	for err != nil {
		evt.Partition = partition
		evt.Topic = TCCConfirmTaskEvent
		evt.Offset = strconv.Itoa(int(offset))
		evt.Retry += 1
		if evt.Retry == MaxTCCMangerProduceRetry {
			//发送到死信队列
			return ErrProduceSubmitFailure
		}
		err = s.TCCMangerProduceAddTCCEvent(evt)
	}
	return nil
} //AddTCCConfirmTask
func NewSaramaSyncProducer(producer sarama.SyncProducer) TCCMegProducer {
	return &SaramaSyncProducer{
		producer: producer,
	}
}
