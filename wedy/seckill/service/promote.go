package service

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/repository"
	"context"
	"errors"
)

type PromoteService interface {
	VerifyCode(ctx context.Context, tccId string, order domain.Order) error
	WithHold(ctx context.Context, activityId int64, productId int64, quality int64) error
	Withdraw(ctx context.Context, productId int64, amount int64) error
}

var (
	ErrPromoteCodeNotFound = errors.New("promote code not found")
	ErrPromoteCodeNotMatch = errors.New("promote code not match")
	ErrPromoteCodeExpired  = errors.New("promote code expired")
)

type promote struct {
	log        logger.LoggerV1
	repo       repository.PromoteCode
	tccManager TccManagerService
}

func (p *promote) VerifyCode(ctx context.Context, tccId string, order domain.Order) error {
	orderTx := domain.OrderTX{
		OrderId:     order.OrderId,
		UserId:      order.UserId,
		Price:       order.Price,
		Quantity:    order.Quantity,
		PromoteCode: order.PromoCode,
	}
	err := p.tccManager.AddTcc(ctx, orderTx, "promote"+tccId)
	if err != nil {
		return err
	}
	err = p.repo.VerifyCode(ctx, tccId, order)
	if err != nil {
		err = p.tccManager.Failed(ctx, "promote"+tccId)
		return err
	}
	err = p.tccManager.Succeed(ctx, "promote"+tccId)
	return err
}

func (p *promote) WithHold(ctx context.Context, activityId int64, productId int64, quality int64) error {
	//TODO implement me
	panic("implement me")
}

func (p *promote) Withdraw(ctx context.Context, productId int64, amount int64) error {
	//TODO implement me
	panic("implement me")
}
