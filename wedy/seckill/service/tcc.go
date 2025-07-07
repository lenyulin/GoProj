package service

import (
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/repository"
	"context"
)

type TccWithDrawnManage interface {
	WithDrawn(ctx context.Context, id string) error
	Failed(ctx context.Context, id string) error
	AddTcc(ctx context.Context, order domain.OrderTX, id string) error
	Succeed(ctx context.Context, id string) error
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
func (t *tcc) Failed(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
func (t *tcc) AddTcc(ctx context.Context, order domain.OrderTX, id string) error {
	return t.repo.AddTcc(ctx, order, id)
}
func (t *tcc) Succeed(ctx context.Context, id string) error {
	return t.repo.TccComplete(ctx, id)
}
