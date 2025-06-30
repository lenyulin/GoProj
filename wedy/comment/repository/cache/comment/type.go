package comment

import (
	"GoProj/wedy/comment/domain"
	"context"
	"errors"
)

type Cache interface {
	//BatchSetAndUpdate(ctx context.Context, key string, comment []domain.Comment) error
	SetAndUpdate(ctx context.Context, key string, comment domain.Comment) error
	Get(ctx context.Context, key string, offset int64) ([]domain.Comment, error)
	IncrLikeCnt(ctx context.Context, id int64, vid int64, uid int64) error
}

var (
	ErrCommentRecordNotFound = errors.New("comment record not found")
	ErrLoadSyncMapComment    = errors.New("failed to load comment")
	ErrCommentCacheExpired   = errors.New("comment cache expired")
)
