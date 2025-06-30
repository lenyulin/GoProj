package dao

import (
	"GoProj/wedy/comment/domain"
	dao2 "GoProj/wedy/internal/repository/dao"
	"context"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GORMCommentDAO struct {
	db *gorm.DB
}

func (d *GORMCommentDAO) IncrLinkeCnt(ctx context.Context, id int64, i int64) error {
	now := time.Now().UnixMilli()
	return d.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"link_cnt": gorm.Expr("link_cnt + ?", 1),
			"utime":    now,
		}),
	}).Create(&Interactive{
		Biz:     biz,
		BizId:   bizId,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}

const (
	MAX_COMMENT_SIZE = 200
)

func (d *GORMCommentDAO) FindById(ctx context.Context, id int64, offset int64) ([]domain.Comment, error) {
	var result []domain.Comment
	err := d.db.WithContext(ctx).
		Where("vid = ?", id).
		Offset(int(offset)).
		Limit(MAX_COMMENT_SIZE).
		Order("ctime DESC").
		Order("like DESC").
		Find(&result).Error
	if err != nil {
		return []domain.Comment{}, err
	}
	return result, nil
}

func NewGORMCommentDAO(db *gorm.DB) CommentDAO {
	return &GORMCommentDAO{
		db: db,
	}
}
func (d *GORMCommentDAO) Insert(ctx context.Context, comment Comment) error {
	now := time.Now().UnixMilli()
	comment.Ctime = now
	err := d.db.WithContext(ctx).Create(&comment).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return dao2.ErrDuplicatedUser
		}
	}
	return err
}
