package repository

import (
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/repository/cache"
	"GoProj/wedy/internal/repository/dao/video"
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type VideoRepository interface {
	Create(ctx context.Context, video domian.Video) error
	Update(ctx context.Context, video domian.Video) error
	SyncStatus(ctx context.Context, video domian.Video) error
	GetByAuthor(ctx context.Context, id int64, limit int, offset int) ([]domian.Video, error)
	GetById(ctx context.Context, id int64) (domian.Video, error)
	ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domian.Video, error)
}
type videoRepository struct {
	dao   dao.VideoDao
	cache cache.VideoCache
}

func NewVideoRepository(dao dao.VideoDao, cache cache.VideoCache) VideoRepository {
	return &videoRepository{
		dao:   dao,
		cache: cache,
	}
}
func (svr *videoRepository) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domian.Video, error) {
	videos, err := svr.dao.ListPub(ctx, start, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.Video, domian.Video](videos, func(idx int, src dao.Video) domian.Video {
		return svr.toDomain(src)
	}), nil
}
func (svr *videoRepository) GetById(ctx context.Context, id int64) (domian.Video, error) {
	result, err := svr.cache.Get(ctx, id)
	if err == nil {
		return result, nil
	}
	video, err := svr.dao.FindById(ctx, id)
	if err != nil {
		return domian.Video{}, err
	}
	go func() {
		err := svr.cache.Set(ctx, svr.toDomain(video))
		if err != nil {
			fmt.Println(err)
		}
	}()
	return svr.toDomain(video), nil
}
func (svr *videoRepository) GetByAuthor(ctx context.Context, id int64, limit int, offset int) ([]domian.Video, error) {
	if offset == 0 && limit == 100 {
		videos, err := svr.cache.GetFirstPage(ctx, id)
		switch err {
		case nil:
			return videos, nil
		default:
			fmt.Println(err)
		}
	}
	videos, err := svr.dao.FindByAuthor(ctx, id, limit, offset)
	if err != nil {
		return nil, err
	}
	res, err := slice.Map[dao.Video, domian.Video](videos, func(idx int, src dao.Video) domian.Video {
		return svr.toDomain(src)
	}), nil
	go func() {
		if offset == 0 && limit == 100 {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = svr.cache.SetFirstPage(ctx, id, res)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = svr.preCache(ctx, videos)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return res, nil
}
func (svr *videoRepository) SyncStatus(ctx context.Context, v domian.Video) error {
	return svr.dao.SyncStatus(ctx, svr.toEntity(v))
}
func (svr *videoRepository) Create(ctx context.Context, v domian.Video) error {
	return svr.dao.Insert(ctx, svr.toEntity(v))
}
func (svr *videoRepository) Update(ctx context.Context, v domian.Video) error {
	return svr.dao.Update(ctx, svr.toEntity(v))
}
func (svr *videoRepository) toDomain(v dao.Video) domian.Video {
	return domian.Video{
		Title:   v.Title,
		Content: v.Content,
		Uid:     v.AuthorId,
		VUid:    v.Id,
	}
}
func (svr *videoRepository) toEntity(v domian.Video) dao.Video {
	return dao.Video{
		Title:    v.Title,
		Content:  v.Content,
		Id:       v.VUid,
		AuthorId: v.Uid,
		Status:   v.Status.ToUint8(),
	}
}

func (svr *videoRepository) preCache(ctx context.Context, videos []dao.Video) error {
	if len(videos) > 0 {
		if err := svr.cache.Set(ctx, svr.toDomain(videos[0])); err != nil {
			return err
		}
	}
	return nil
}
