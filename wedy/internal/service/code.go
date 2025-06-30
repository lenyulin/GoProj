package service

import (
	"GoProj/wedy/internal/repository"
	"GoProj/wedy/internal/service/sms"
	"context"
	"errors"
	"fmt"
	"math/rand"
)

var (
	ErrCodeVerifySendTooMany = repository.ErrCodeVerifySendTooMany
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}
type codeService struct {
	repo      repository.CodeRepository
	sms       sms.Service
	codeTplId string
}

func NewCodeService(repo repository.CodeRepository, sms sms.Service) CodeService {
	return &codeService{
		repo:      repo,
		sms:       sms,
		codeTplId: "codeTplId",
	}
}

func (svc *codeService) generateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	return svc.sms.Send(ctx, svc.codeTplId, []string{code}, phone)
}

func (svc *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if errors.Is(err, repository.ErrCodeVerifySendTooMany) {
		return false, err
	}
	return ok, err
}
