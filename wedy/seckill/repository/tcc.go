package repository

import (
	"GoProj/wedy/seckill/repository/dao"
	"context"
)

type Tcc interface {
	CancelTcc(ctx context.Context, id string) error
	AddTcc(ctx context.Context, id string) error
	TccComplete(ctx context.Context, id string) error
	TccFailed(ctx context.Context, id string) error
}

type tcc struct {
	dao dao.Tcc
}

func NewTccRepository(dao dao.Tcc) Tcc {
	return &tcc{
		dao: dao,
	}
}

func (t *tcc) CancelTcc(ctx context.Context, id string) error {
	return t.dao.Cancel(ctx, id)
}

func (t *tcc) AddTcc(ctx context.Context, id string) error {
	return t.dao.Add(ctx, id)
}

func (t *tcc) TccComplete(ctx context.Context, id string) error {
	return t.dao.Complete(ctx, id)
}

func (t *tcc) TccFailed(ctx context.Context, id string) error {
	return t.dao.Failed(ctx, id)
}
