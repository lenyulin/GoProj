package dao

import (
	"GoProj/wedy/internal/repository/cache"
	"GoProj/wedy/internal/repository/dao"
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

type GORMVideoDAO struct {
	db    *gorm.DB
	cache cache.VideoCache
}

func NewGORMVideoDAO(db *gorm.DB, cache cache.VideoCache) *GORMVideoDAO {
	return &GORMVideoDAO{
		db:    db,
		cache: cache,
	}
}

const VideoStatusNotPublished = 2

func (d *GORMVideoDAO) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]Video, error) {
	var videos []Video
	err := d.db.WithContext(ctx).Where("utime < ? AND status = ?", start.UnixMilli(), VideoStatusNotPublished).
		Offset(offset).
		Limit(limit).
		Order("utime DESC").
		Find(&videos).
		Error
	return videos, err
}
func (d *GORMVideoDAO) FindByAuthor(ctx context.Context, id int64, limit int, offset int) ([]Video, error) {
	var videos []Video
	err := d.db.WithContext(ctx).Where("author_id = ?", id).Offset(offset).Limit(limit).Order("utime DESC").Find(&videos).Error
	return videos, err
}
func (d *GORMVideoDAO) FindById(ctx context.Context, id int64) (Video, error) {
	_, err := d.cache.Get(ctx, id)
	if err == nil {
		return Video{}, err
	}
	var video Video
	err = d.db.WithContext(ctx).Where("id = ?", id).First(&video).Error
	return video, err
}
func (d *GORMVideoDAO) Update(ctx context.Context, v Video) error {
	v.Utime = time.Now().Unix()
	res := d.db.WithContext(ctx).Model(&v).
		Where("V_Uid = ? AND Author_Id = ?", v.Id, v.AuthorId).
		Updates(map[string]any{
			"Title":   v.Title,
			"Content": v.Content,
			"Status":  v.Status,
			"Utime":   v.Utime,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("ID or AuthorId does not exist")
	}
	return nil
}
func (d *GORMVideoDAO) SyncStatus(ctx context.Context, v Video) error {
	now := time.Now().UnixMilli()
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Video{}).Where("Id = ? AND Author_Id = ?", v.Id, v.AuthorId).
			Updates(map[string]any{
				"Status": v.Status,
				"Utime":  now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errors.New("Error unauthorized access")
		}
		return nil
	})
}

func (d *GORMVideoDAO) Insert(ctx context.Context, v Video) error {
	now := time.Now().UnixMilli()
	v.Ctime = now
	v.Utime = now
	err := d.db.WithContext(ctx).Create(&v).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return dao.ErrDuplicatedUser
		}
	}
	return err
}
