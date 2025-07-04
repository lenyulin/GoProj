package repository

import (
	"GoProj/wedy/seckill/repository/dao"
	"context"
)

type Tcc interface {
	CancelTcc(ctx context.Context, id string) error
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
	return t.dao.SubmitCancelRequest(ctx, id)
}
