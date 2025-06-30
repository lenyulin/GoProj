package service

import (
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/repository"
	repomock "GoProj/wedy/internal/repository/mock"
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/assert/v2"
	assert2 "github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPWDEncrypt(t *testing.T) {
	pwd := []byte("Hello123456")
	encryptedPWD, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	assert2.NoError(t, err)
	fmt.Println(string(encryptedPWD))
}
func TestUserSignupSVC(t *testing.T) {
	testcases := []struct {
		name     string
		wantErr  error
		ctx      context.Context
		Phone    string
		Password string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
	}{
		{
			name:     "success",
			wantErr:  nil,
			Phone:    "13333333333",
			Password: "Hello123456",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return userRepo
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSVC := NewUserService(tc.mock(ctrl))
			err := userSVC.Signup(context.Background(), domian.User{
				Phone:    tc.Phone,
				Password: tc.Password,
			})
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
func TestUserLoginSVC(t *testing.T) {
	testcases := []struct {
		name           string
		wantErr        error
		ctx            context.Context
		Phone          string
		Password       string
		wantDomainUser domian.User
		mock           func(ctrl *gomock.Controller) repository.UserRepository
	}{
		{
			name:     "success",
			wantErr:  nil,
			Phone:    "13333333333",
			Password: "Hello123456",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByPhone(context.Background(), "13333333333").Return(domian.User{
					Phone:    "13333333333",
					Password: "$2a$10$Qig1Lz1/fd2KsKPQ.E7lquwzexFqzGNJk/FmTvHjzBJUnkfw9Ntle",
				}, nil)
				return userRepo
			},
			wantDomainUser: domian.User{
				Phone: "13333333333",
			},
		},
		{
			name:     "ErrInvalidUserOrPassword",
			wantErr:  ErrInvalidUserOrPassword,
			Phone:    "13333333333",
			Password: "Hello12345",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByPhone(context.Background(), "13333333333").Return(domian.User{
					Phone:    "13333333333",
					Password: "$2a$10$Qig1Lz1/fd2KsKPQ.E7lquwzexFqzGNJk/FmTvHjzBJUnkfw9Ntle",
				}, nil)
				return userRepo
			},
			wantDomainUser: domian.User{},
		},
		{
			name:     "UserNOTFound",
			wantErr:  ErrInvalidUserOrPassword,
			Phone:    "13333333333",
			Password: "Hello12345",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByPhone(context.Background(), "13333333333").Return(domian.User{}, repository.ErrUserNotFound)
				return userRepo
			},
			wantDomainUser: domian.User{},
		},
		{
			name:     "OthersErrors",
			wantErr:  errors.New("some error"),
			Phone:    "13333333333",
			Password: "Hello12345",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByPhone(context.Background(), "13333333333").Return(domian.User{}, errors.New("some error"))
				return userRepo
			},
			wantDomainUser: domian.User{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSVC := NewUserService(tc.mock(ctrl))
			u, err := userSVC.Login(context.Background(), tc.Phone, tc.Password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantDomainUser, u)
		})
	}
}
func TestUserFindOrCreateSVC(t *testing.T) {
	testcases := []struct {
		name           string
		wantErr        error
		ctx            context.Context
		Phone          string
		Password       string
		wantDomainUser domian.User
		mock           func(ctrl *gomock.Controller) repository.UserRepository
	}{
		{
			name:     "success",
			wantErr:  nil,
			Phone:    "13333333333",
			Password: "Hello123456",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByPhone(context.Background(), "13333333333").Return(domian.User{}, repository.ErrUserNotFound)
				userRepo.EXPECT().Create(context.Background(), domian.User{
					Phone: "13333333333",
				}).Return(nil)
				userRepo.EXPECT().FindByPhone(context.Background(), "13333333333").Return(domian.User{
					Phone: "13333333333",
				}, nil)
				return userRepo
			},
			wantDomainUser: domian.User{
				Phone: "13333333333",
			},
		},
		{
			name:     "UserExits",
			wantErr:  nil,
			Phone:    "13333333333",
			Password: "Hello123456",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByPhone(context.Background(), "13333333333").Return(domian.User{
					Phone: "13333333333",
				}, nil)
				return userRepo
			},
			wantDomainUser: domian.User{
				Phone: "13333333333",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSVC := NewUserService(tc.mock(ctrl))
			u, err := userSVC.FindOrCreate(context.Background(), tc.Phone)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantDomainUser, u)
		})
	}
}
