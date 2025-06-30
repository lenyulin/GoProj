package repository

import (
	"GoProj/wedy/comment/domain"
	"errors"
)

type Consumer interface {
	Start() error
}

const (
	TopicCommentSubmitEvent      = "comment_submit"
	TopicCommentIncrLikeCntEvent = "comment_incr_like"
)

var (
	ErrUpdatePrimaryIncrLikeCnt = errors.New("update local cache error")
	ErrProducerIncrLikeCnt      = errors.New("producer incr like cnt error")
)

type CommentEvent struct {
	//Comment ID
	Id        int64
	partition []string
	offset    []string
	timeStamp []int64
	topic     string
	retry     int64
	comment   domain.Comment
}

type CommentLikeEvent struct {
	//Comment ID
	Id        int64
	Uid       int64
	Vid       int64
	partition []string
	offset    []string
	timeStamp []int64
	topic     string
	retry     int64
}
