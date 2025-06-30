package domian

import "time"

type Video struct {
	Title   string
	Content string
	Uid     int64
	VUid    int64
	Status  VideoStatus
	Author  Author
	CTime   time.Time
	Utime   time.Time
}
type VideoStatus uint8

func (v VideoStatus) ToUint8() uint8 {
	return uint8(v)
}

const (
	VideoStatusUnknown = iota
	VideoStatusUnpublished
	VideoStatusPublished
	VideoStatusPrivate
)

type Author struct {
	Name string
	Id   int64
}
