package service

import (
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/repository"
	"context"
	"errors"
)

type PromoteService interface {
	VerifyCode(ctx context.Context, order domain.Order) error
	WithHold(ctx context.Context, activityId int64, productId int64, quality int64) error
	Withdraw(ctx context.Context, productId int64, amount int64) error
}

var (
	ErrPromoteCodeNotFound = errors.New("promote code not found")
	ErrPromoteCodeNotMatch = errors.New("promote code not match")
	ErrPromoteCodeExpired  = errors.New("promote code expired")
)

type promote struct {
	log  logger.LoggerV1
	repo repository.PromoteCode
}

func (p *promote) VerifyCode(ctx context.Context, activityId int64, productId int64, code []string) error {
	//TODO implement me
	panic("implement me")
}

func (p *promote) WithHold(ctx context.Context, activityId int64, productId int64, quality int64) error {
	//TODO implement me
	panic("implement me")
}

func (p *promote) Withdraw(ctx context.Context, productId int64, amount int64) error {
	//TODO implement me
	panic("implement me")
}
