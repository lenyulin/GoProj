package service

import (
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicatedUser        = repository.ErrDuplicatedUser
	ErrInvalidUserOrPassword = errors.New("invalid user or password")
)

type UserService interface {
	Signup(ctx context.Context, u domian.User) error
	Login(ctx context.Context, phone string, password string) (domian.User, error)
	FindOrCreate(ctx context.Context, phone string) (domian.User, error)
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx context.Context, u domian.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, phone string, password string) (domian.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domian.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domian.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domian.User{}, ErrInvalidUserOrPassword
	}
	return domian.User{
		Phone: u.Phone,
	}, nil
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domian.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return u, err
	}
	err = svc.repo.Create(ctx, domian.User{
		Phone: phone,
	})
	// 1.err唯一索引错误(phone存在) 2.err!=nil interal error
	if err != nil && !errors.Is(err, repository.ErrDuplicatedUser) {
		return domian.User{}, err
	}
	// err==nil or ErrDuplicatedUser
	// 主从延迟
	return svc.repo.FindByPhone(ctx, phone)
}
