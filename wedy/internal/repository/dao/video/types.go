package dao

import (
	"context"
	"time"
)

type VideoDao interface {
	Insert(ctx context.Context, video Video) error
	Update(ctx context.Context, video Video) error
	SyncStatus(ctx context.Context, video Video) error
	FindByAuthor(ctx context.Context, id int64, limit int, offset int) ([]Video, error)
	FindById(ctx context.Context, id int64) (Video, error)
	ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]Video, error)
}

type Video struct {
	Title    string `gorm:"varchar(255)" bson:"title"`
	Content  string `gorm:"varchar(4096)" bson:"content"`
	Id       int64  `gorm:"primaryKey" bson:"id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	AuthorId int64  `gorm:"index" bson:"authorId,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
