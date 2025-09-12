package tcc

import (
	"GoProj/wedy/seckill/service"
	"context"
)

type OrderAdapter struct {
	orderService service.OrderService // 依赖库存服务
}

func (o *OrderAdapter) Try(ctx context.Context, tccID string, bizData interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (o *OrderAdapter) Confirm(ctx context.Context, tccID string) error {
	//TODO implement me
	panic("implement me")
}

func (o *OrderAdapter) Cancel(ctx context.Context, tccID string) error {
	//TODO implement me
	panic("implement me")
}
