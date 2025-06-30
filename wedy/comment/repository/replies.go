package repository

import (
	"GoProj/wedy/comment/domain"
	"GoProj/wedy/comment/repository/cache/replies"
	"context"
	"sync"
)

type ReplyRepository interface {
	Submit(ctx context.Context, comment domain.Comment) error
	Get(ctx context.Context, id int64, page int64) ([]domain.Comment, error)
}
type CompositeReplyCache struct {
	primary replies.Cache // 主缓存（如Redis）
	//secondary comment2.CommentCache // 二级缓存（如本地）
}

type replyRepository struct {
	replyCache CompositeReplyCache
	mutex      sync.RWMutex
}

func (r *replyRepository) Submit(ctx context.Context, comment domain.Comment) error {
	//TODO implement me
	panic("implement me")
}

func (r *replyRepository) Get(ctx context.Context, id int64, page int64) ([]domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func NewReplyRepository(primary replies.Cache) ReplyRepository {
	return &replyRepository{
		replyCache: CompositeReplyCache{
			primary: primary,
		},
	}
}
