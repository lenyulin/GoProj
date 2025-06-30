package dao

import (
	"GoProj/wedy/internal/domian"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type JobDAO interface {
	Preempt(ctx context.Context) (domian.Job, error)
	Release(ctx context.Context, id int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
}

const (
	jobStatusWaiting = iota
	jobStatusRunning
	jobStatusPaused
)

type Job struct {
	Id         int64 `gorm:"primary_key;auto_increment"`
	Status     int
	Version    int
	Utime      int64
	Ctime      int64
	NextTime   int64 `gorm:"index"`
	Executor   string
	Name       string `gorm:"unique;type:varchar(128)"`
	Expression string
}

type GROMJobDAO struct {
	db *gorm.DB
}

func (dao *GROMJobDAO) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ? ?", id).Updates(map[string]any{
		"utime":     now,
		"next_time": t.UnixMilli(),
	}).Error
}
func (dao *GROMJobDAO) UpdateUtime(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ? ?", id).Updates(map[string]any{
		"utime": now,
	}).Error
}
func (dao *GROMJobDAO) Preempt(ctx context.Context) (domian.Job, error) {
	db := dao.db.WithContext(ctx)
	for {
		var j Job
		now := time.Now().UnixMilli()
		err := db.
			Where("WHERE status = ? AND next_time<", jobStatusWaiting, now).
			First(&j.Id).
			Error
		if err != nil {
			return domian.Job{}, err
		}
		res := db.Model(&j).Where("id = ? AND version = ?", j.Id, j.Version).Updates(map[string]any{
			"status":  jobStatusRunning,
			"version": j.Version + 1,
			"utime":   now,
		})
		if res.Error != nil {
			return domian.Job{}, res.Error
		}
		if res.RowsAffected == 0 {
			return domian.Job{}, errors.New("preempt job failed")
		}
		return domian.Job{
			Id: j.Id,
		}, err
	}
}

func (dao *GROMJobDAO) Release(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ? ?", id).Updates(map[string]any{
		"status": jobStatusPaused,
		"utime":  now,
	}).Error
}
