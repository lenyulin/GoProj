package dao

import (
	"GoProj/wedy/pkg/oss"
	"context"
	"errors"
	"time"
)

var (
	ErrRecordInsertFailed = errors.New("Insert Video Info to DB error")
	ErrVideoUploadFailed  = errors.New("Upload Video error")
)

type OSSVideoDao struct {
	OSS     oss.OSSHandler
	gormDAO VideoDao
}

func (d *OSSVideoDao) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]Video, error) {
	//TODO implement me
	panic("implement me")
}

func (d *OSSVideoDao) FindByAuthor(ctx context.Context, id int64, limit int, offset int) ([]Video, error) {
	//TODO implement me
	panic("implement me")
}

func (d *OSSVideoDao) FindById(ctx context.Context, id int64) (Video, error) {
	//TODO implement me
	panic("implement me")
}

func NewOSSVideoDao(OSS oss.OSSHandler, gormDAO VideoDao) *OSSVideoDao {
	return &OSSVideoDao{OSS: OSS, gormDAO: gormDAO}
}
func (d *OSSVideoDao) Update(ctx context.Context, v Video) error {
	return d.gormDAO.Update(ctx, v)
}

func (d *OSSVideoDao) Insert(ctx context.Context, video Video) error {
	err := d.gormDAO.Insert(ctx, video)
	if err != nil {
		return ErrRecordInsertFailed
	}
	if err = d.OSS.Upload(ctx, video.Id); err != nil {
		return ErrVideoUploadFailed
	}
	return nil
}

func (d *OSSVideoDao) SyncStatus(ctx context.Context, v Video) error {
	return d.gormDAO.SyncStatus(ctx, v)
}

func (d *OSSVideoDao) insert(ctx context.Context, v Video) error {
	return d.gormDAO.Insert(ctx, v)
}
