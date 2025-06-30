package dao

import (
	"GoProj/wedy/comment/domain"
	"context"
)

type Comment struct {
	Id             int64           `json:"comment_id"`
	Uid            int64           `json:"user_id"`
	Content        string          `json:"content"`
	Ctime          int64           `json:"ctime"`
	Like           int64           `json:"like"`
	Dislike        int64           `json:"dislike"`
	Score          int64           `json:"score"`
	ReplyCount     int64           `json:"reply_count"`
	Picture        string          `json:"picture"`
	MentionedUsers []MentionedUser `json:"mentioned_users"`
}

type MentionedUser struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type CommentDAO interface {
	Insert(ctx context.Context, comment Comment) error
	FindById(ctx context.Context, id int64, offset int64) ([]domain.Comment, error)
	IncrLinkeCnt(ctx context.Context, id int64, i int64) error
}
