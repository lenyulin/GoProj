package service

import (
	"GoProj/wedy/comment/domain"
	"GoProj/wedy/comment/repository"
	"GoProj/wedy/pkg/snowflake"
	"context"
)

type CommentService interface {
	Submit(ctx context.Context, comment domain.Comment) error
	Get(ctx context.Context, id int64, page int64) ([]domain.Comment, error)
	Like(ctx context.Context, id int64, vid int64, uid int64) error
}
type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{
		repo: repo,
	}
}
func (svc *commentService) Like(ctx context.Context, id int64, vid int64, uid int64) error {
	return svc.repo.IncrLikeCnt(ctx, id, vid, uid)
}
func (svc *commentService) Get(ctx context.Context, id int64, page int64) ([]domain.Comment, error) {
	return svc.repo.Get(ctx, id, page)
}
func (svc *commentService) Submit(ctx context.Context, comment domain.Comment) error {
	uid, err := snowflake.Generate()
	if err != nil {
		return ErrGenerateSnowFlakeError
	}
	comment.Id = uid
	err = svc.repo.Submit(ctx, comment)
	if err != nil {
		return err
	}
	return nil
}
