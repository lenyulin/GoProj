package dao

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	ErrBucketOpenUploadStreamFaild = errors.New("ErrBucketOpenUploadStreamFaild")
	ErrBucketUploadStreamTimeOut   = errors.New("ErrBucketUploadStreamTimeOut")
	ErrOpenFileFaild               = errors.New("ErrOpenFileFaild")
)

type MongoDBDAO struct {
	gorm   *GORMVideoDAO
	mdb    *mongo.Database
	col    *mongo.Collection
	bucket *gridfs.Bucket
}

func NewMongoDBDAO(mdb *mongo.Database, col *mongo.Collection, bucket *gridfs.Bucket, gorm *GORMVideoDAO) *MongoDBDAO {
	return &MongoDBDAO{
		mdb:    mdb,
		col:    col,
		bucket: bucket,
		gorm:   gorm,
	}
}

const Dst = "C:\\Users\\lyl69\\GolandProjects\\GoProj\\wedy\\tmp\\"

func (m MongoDBDAO) Insert(ctx context.Context, v Video) error {
	v.Ctime = time.Now().UnixMilli()
	v.Utime = time.Now().UnixMilli()
	metadata := bson.M{
		"Title":    v.Title,
		"Content":  v.Content,
		"AuthorId": v.AuthorId,
		"id":       v.Id,
		"Status":   v.Status,
		"Ctime":    v.Ctime,
		"Utime":    v.Utime,
	}
	// 打开文件
	file, err := os.Open(fmt.Sprintf("%s%s%s", Dst, strconv.FormatInt(v.Id, 10), ".mp4"))
	if err != nil {
		fmt.Printf("cannot open file: %v\n", err)
		return ErrOpenFileFaild
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return ErrOpenFileFaild
	}
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)
	_, err = file.Read(buffer)
	if err != nil {
		return ErrOpenFileFaild
	}
	uploadOpts := options.GridFSUpload().
		SetMetadata(metadata)
	uploadStream, err := m.bucket.OpenUploadStream(
		fmt.Sprintf(strconv.FormatInt(v.Id, 10)),
		uploadOpts,
	)
	if err != nil {
		return ErrBucketOpenUploadStreamFaild
	}
	defer func() {
		if err = uploadStream.Close(); err != nil {
			log.Panic(err)
		}
	}()
	err = uploadStream.SetWriteDeadline(time.Now().Add(30 * time.Second))
	if err != nil {
		return ErrBucketUploadStreamTimeOut
	}
	if _, err = uploadStream.Write(buffer); err != nil {
		log.Panic(err)
	}
	return nil
}

func (m MongoDBDAO) Update(ctx context.Context, v Video) error {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBDAO) SyncStatus(ctx context.Context, v Video) error {
	//TODO implement me
	panic("implement me")
}
