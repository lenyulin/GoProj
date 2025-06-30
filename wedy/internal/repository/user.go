package repository

import (
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/repository/cache"
	"GoProj/wedy/internal/repository/dao"
	"context"
	"database/sql"
	"fmt"
)

var (
	ErrDuplicatedUser = dao.ErrDuplicatedUser
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, user domian.User) error
	FindByPhone(ctx context.Context, phone string) (domian.User, error)
}
type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCachedUserRepository(dao dao.UserDAO, uc cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: uc,
	}
}
func (repo *CachedUserRepository) Create(ctx context.Context, user domian.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(user))
}
func (repo *CachedUserRepository) toEntity(user domian.User) dao.User {
	return dao.User{
		Id:       user.Id,
		Phone:    sql.NullString{String: user.Phone, Valid: user.Phone != ""},
		Password: user.Password,
	}
}
func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domian.User, error) {
	du, err := repo.cache.Get(ctx, phone)
	switch err {
	case nil:
		return du, nil
	case cache.ErrKeyNotExists:
		u, err := repo.dao.FindByPhone(ctx, phone)
		if err != nil {
			return domian.User{}, err
		}
		du = repo.toDomainUser(u)
		go func() {
			err = repo.cache.Set(ctx, du)
			if err != nil {
				fmt.Println(err)
			}
		}()
		return repo.toDomainUser(u), nil
	default:
		return domian.User{}, err
	}
}

func (repo *CachedUserRepository) toDomainUser(u dao.User) domian.User {
	return domian.User{
		Phone:    u.Phone.String,
		Password: u.Password,
		Id:       u.Id,
	}
}
