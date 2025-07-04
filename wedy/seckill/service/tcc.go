package service

import (
	"GoProj/wedy/seckill/repository"
	"context"
)

type TccWithDrawnManage interface {
	WithDrawn(ctx context.Context, id string) error
	CheckTcc(ctx context.Context, id string) error
}
type tcc struct {
	repo repository.Tcc
}

func NewTccService(repo repository.Tcc) TccWithDrawnManage {
	return &tcc{
		repo: repo,
	}
}

func (t *tcc) WithDrawn(ctx context.Context, id string) error {
	return t.repo.CancelTcc(ctx, id)
}
func (t *tcc) CheckTcc(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
