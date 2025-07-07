package repository

import (
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/repository/dao"
	"context"
)

type Tcc interface {
	CancelTcc(ctx context.Context, order domain.OrderTX, id string) error
	AddTcc(ctx context.Context, order domain.OrderTX, id string) error
	TccComplete(ctx context.Context, order domain.OrderTX, id string) error
	TccFailed(ctx context.Context, order domain.OrderTX, id string) error
}

type tcc struct {
	dao dao.Tcc
}

func NewTccRepository(dao dao.Tcc) Tcc {
	return &tcc{
		dao: dao,
	}
}
func (t *tcc) toTxDAO(src domain.OrderTX) dao.OrderTXDAO {
	return dao.OrderTXDAO{
		OrderId:     src.OrderId,
		UserId:      src.UserId,
		Price:       src.Price,
		Quantity:    src.Quantity,
		PromoteCode: src.PromoteCode,
	}
}
func (t *tcc) CancelTcc(ctx context.Context, order domain.OrderTX, id string) error {
	return t.dao.Cancel(ctx, t.toTxDAO(order), id)
}

func (t *tcc) AddTcc(ctx context.Context, order domain.OrderTX, id string) error {
	return t.dao.Add(ctx, t.toTxDAO(order), id)
}

func (t *tcc) TccComplete(ctx context.Context, order domain.OrderTX, id string) error {
	return t.dao.Complete(ctx, t.toTxDAO(order), id)
}

func (t *tcc) TccFailed(ctx context.Context, order domain.OrderTX, id string) error {
	return t.dao.Failed(ctx, t.toTxDAO(order), id)
}
