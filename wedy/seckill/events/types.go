package events

type Consumer interface {
	Start() error
}

type AddTCCEvent struct {
	TCCIdx    int64
	Partition int32
	Offset    string
	TimeStamp int64
	Topic     string
	Retry     int64
	DATA      interface{}
}

const (
	TopicWatchEvent     = "tcc_manger_watch"
	TCCConfirmTaskEvent = "tcc_confirm_task"
)
