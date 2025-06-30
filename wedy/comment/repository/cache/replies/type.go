package replies

import (
	"GoProj/wedy/comment/domain"
	"context"
	"errors"
)

type Cache interface {
	SetAndUpdate(ctx context.Context, key string, comment domain.Reply) error
	Get(ctx context.Context, key string, offset int64) ([]domain.Reply, error)
}

var (
	ErrReplyRecordNotFound = errors.New("reply record not found")
	ErrLoadAtomicReply     = errors.New("failed to load reply")
	ErrReplyCacheExpired   = errors.New("reply cache expired")
)
