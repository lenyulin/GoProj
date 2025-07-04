package event

import "GoProj/wedy/comment/domain"

type Consumer interface {
	Start() error
}

const TccCanceledEvent = "tcc_canceled"

type TccCancelEvent struct {
	//Comment ID
	Id        int64
	partition []string
	offset    []string
	timeStamp []int64
	topic     string
	retry     int64
	comment   domain.Comment
}
