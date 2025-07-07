package service

import (
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/repository"
	"context"
	"encoding/json"
	"errors"
	"github.com/allegro/bigcache/v3"
	"time"
)

type TccManagerService interface {
	WithDrawn(ctx context.Context, id string) error
	Failed(ctx context.Context, id string) error
	AddTcc(ctx context.Context, order domain.OrderTX, id string) error
	Succeed(ctx context.Context, id string) error
}
type tcc struct {
	repo              repository.Tcc
	localMessageTable *bigcache.BigCache
}

var (
	ErrOrderTxExist = errors.New("order tx exist")
)

const (
	OrderTxCreated = iota + 1
)

func NewTccService(repo repository.Tcc, localMessageTable *bigcache.BigCache) TccManagerService {
	return &tcc{
		repo:              repo,
		localMessageTable: localMessageTable,
	}
}

func (t *tcc) WithDrawn(ctx context.Context, id string) error {
	//return t.repo.CancelTcc(ctx, id)
	//TODO implement me
	panic("implement me")
}
func (t *tcc) Failed(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
func (t *tcc) AddTcc(ctx context.Context, order domain.OrderTX, id string) error {
	_, err := t.localMessageTable.Get(id)
	if err != nil {
		return ErrOrderTxExist
	}
	localTx := domain.LocalTxMessage{
		Status:        OrderTxCreated,
		RetryCount:    0,
		CTime:         time.Now().UnixMilli(),
		UTime:         time.Now().UnixMilli(),
		NextRetryTime: 1,
		Tx:            order,
	}
	lTx, err := json.Marshal(localTx)
	if err != nil {
		return err
	}
	err = t.localMessageTable.Set("id", lTx)
	if err != nil {
		return err
	}
	err = t.repo.AddTcc(ctx, order, id)
	return err
}
func (t *tcc) Succeed(ctx context.Context, id string) error {
	//return t.repo.TccComplete(ctx, id)
	//TODO implement me
	panic("implement me")
}
