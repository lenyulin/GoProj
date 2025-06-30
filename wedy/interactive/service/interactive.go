package service

import (
	"GoProj/wedy/interactive/domain"
	"GoProj/wedy/interactive/repository"
	"context"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error)
}
type interactiveService struct {
	repo repository.InteractiveRepository
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo: repo,
	}
}

func (s *interactiveService) GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error) {
	intres, err := s.repo.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.Interactive)
	for _, intr := range intres {
		res[intr.BizId] = intr
	}
	return res, nil
}
func (s *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return s.repo.IncrReadCnt(ctx, biz, bizId)
}
