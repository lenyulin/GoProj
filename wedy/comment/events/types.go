package events

import "GoProj/wedy/comment/domain"

type Consumer interface {
	Start() error
}

const TopicCommentSubmitEvent = "comment_submit"

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
