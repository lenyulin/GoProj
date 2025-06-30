package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"strconv"
	"time"
)

type CommentProducer interface {
	CommentProduceSubmitEvent(evt CommentEvent) error
	CommentIncrLikeCntEvent(evt CommentLikeEvent) error
}

const (
	MaxCommentProduceRetry     = 3
	MaxCommentLikeProduceRetry = 3
)

var (
	ErrProduceSubmitFailure = errors.New("failed to produce sumit event")
)

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func (s *SaramaSyncProducer) CommentIncrLikeCntEvent(evt CommentLikeEvent) error {
	val, err := json.Marshal(&evt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	evt.timeStamp = append(evt.timeStamp, time.Now().UnixMilli())
	partition, offset, err := s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicCommentIncrLikeCntEvent,
		Value: sarama.StringEncoder(val),
	})
	if err != nil {
		evt.partition = append(evt.partition, strconv.Itoa(int(partition)))
		evt.topic = TopicCommentIncrLikeCntEvent
		evt.offset = append(evt.offset, strconv.Itoa(int(offset)))
		evt.retry += 1
		if evt.retry == MaxCommentLikeProduceRetry {
			//发送到死信队列
			return ErrProduceSubmitFailure
		}
		err = s.CommentIncrLikeCntEvent(evt)
		return err
	}
	return nil
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) CommentProducer {
	return &SaramaSyncProducer{
		producer: producer,
	}
}

func (s *SaramaSyncProducer) CommentProduceSubmitEvent(evt CommentEvent) error {
	val, err := json.Marshal(&evt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	evt.timeStamp = append(evt.timeStamp, time.Now().UnixMilli())
	partition, offset, err := s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicCommentSubmitEvent,
		Value: sarama.StringEncoder(val),
	})
	if err != nil {
		evt.partition = append(evt.partition, strconv.Itoa(int(partition)))
		evt.topic = TopicCommentSubmitEvent
		evt.offset = append(evt.offset, strconv.Itoa(int(offset)))
		evt.retry += 1
		if evt.retry == MaxCommentProduceRetry {
			//发送到死信队列
			return ErrProduceSubmitFailure
		}
		err = s.CommentProduceSubmitEvent(evt)
		return err
	}
	return nil
}
