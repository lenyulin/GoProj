package repository

import "context"

type Inventory interface {
	ReserveStock(ctx context.Context, activityId int64, productId int64, quantity int64) error
	Withdraw(ctx context.Context, productId int64, amount int64) error
}
type inventoryRepo struct{}

func (i *inventoryRepo) ReserveStock(ctx context.Context, activityId int64, productId int64, quantity int64) error {
	//TODO implement me
	panic("implement me")
}

func (i *inventoryRepo) Withdraw(ctx context.Context, productId int64, amount int64) error {
	//TODO implement me
	panic("implement me")
}
