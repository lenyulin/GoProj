package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicatedUser = errors.New("duplicated user")
	ErrRecordNotFound = errors.New("user not found")
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByPhone(ctx context.Context, phone string) (User, error)
}
type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return ErrDuplicatedUser
		}
	}
	return err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var usr User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&usr).Error
	return usr, err
}

type User struct {
	Id       int64          `gorm:"primary_key;auto_increment;unique"`
	Phone    sql.NullString `gorm:"unique"`
	Password string
	Ctime    int64
	Utime    int64
}
