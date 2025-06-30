package repository

import (
	"GoProj/wedy/comment/domain"
	"GoProj/wedy/comment/repository/cache/comment"
	"GoProj/wedy/comment/repository/dao"
	"GoProj/wedy/pkg/logger"
	"context"
	"errors"
	"github.com/ecodeclub/ekit/slice"
	"strconv"
)

type CommentRepository interface {
	Submit(ctx context.Context, comment domain.Comment) error
	Get(ctx context.Context, id int64, page int64) ([]domain.Comment, error)
	IncrLikeCnt(ctx context.Context, id int64, vid int64, uid int64) error
}
type CompositeCommentCache struct {
	primary   comment.Cache // 本地主缓存
	secondary comment.Cache // 二级redis缓存
}

type commentRepository struct {
	commentCache CompositeCommentCache
	commentDAO   dao.CommentDAO
	producer     CommentProducer
	log          logger.LoggerV1
}

func NewCommentRepository(primary comment.Cache, secondary comment.Cache, commentDAO dao.CommentDAO, producer CommentProducer, log logger.LoggerV1) CommentRepository {
	return &commentRepository{
		commentCache: CompositeCommentCache{
			primary:   primary,
			secondary: secondary,
		},
		commentDAO: commentDAO,
		producer:   producer,
		log:        log,
	}
}

func (c *commentRepository) IncrLikeCnt(ctx context.Context, id int64, vid int64, uid int64) error {
	//更新本地一级缓存
	err := c.commentCache.primary.IncrLikeCnt(ctx, id, vid, uid)
	if err != nil {
		return ErrUpdatePrimaryIncrLikeCnt
	}
	//更新redis次级缓存
	go func() {
		er := c.commentCache.secondary.IncrLikeCnt(ctx, id, vid, uid)
		if er != nil {
			c.log.Warn("ErrUpdateSecondaryIncrLikeCnt")
		}
	}()

	//发送到kafka批量更新 ctx, id, vid, uid
	err = c.producer.CommentIncrLikeCntEvent(CommentLikeEvent{
		Id:  id,
		Uid: uid,
		Vid: vid,
	})
	if err != nil {
		return ErrProducerIncrLikeCnt
	}
	return nil
}

func (c *commentRepository) Submit(ctx context.Context, comment domain.Comment) error {
	err := c.commentCache.primary.SetAndUpdate(ctx, strconv.FormatInt(comment.VId, 10), comment)
	if err != nil {
		return err
	}
	go func() {
		c.commentCache.secondary.SetAndUpdate(ctx, strconv.FormatInt(comment.VId, 10), comment)
	}()
	err = c.commentDAO.Insert(ctx, c.toEntity(comment))
	if err != nil {
		return err
	}
	return nil
}

func (c *commentRepository) toEntity(comment domain.Comment) dao.Comment {
	users := slice.Map[domain.MentionedUser, dao.MentionedUser](comment.MentionedUsers, func(idx int, src domain.MentionedUser) dao.MentionedUser {
		return dao.MentionedUser{
			Id:   src.Id,
			Name: src.Name,
		}
	})
	return dao.Comment{
		Id:             comment.Id,
		Uid:            comment.User.Id,
		Ctime:          comment.Ctime,
		Content:        comment.Content,
		Like:           comment.Like,
		Dislike:        comment.Dislike,
		ReplyCount:     comment.ReplyCount,
		Picture:        comment.Picture,
		MentionedUsers: users,
	}
}

func (c *commentRepository) Get(ctx context.Context, id int64, page int64) ([]domain.Comment, error) {
	res, err := c.commentCache.primary.Get(ctx, strconv.FormatInt(id, 10), page*10)
	if err == nil {
		return res, nil
	}
	res, err = c.commentCache.secondary.Get(ctx, strconv.FormatInt(id, 10), page*10)
	if err == nil {
		return res, nil
	}
	if errors.Is(err, comment.ErrCommentRecordNotFound) {
		//查询数据库
		res, err = c.commentDAO.FindById(ctx, id, page)
		if err != nil {
			return nil, err
		}
		//更新缓存
		err = c.commentCache.primary.SetAndUpdate(ctx, strconv.FormatInt(id, 10), res[])
		go func() {
			c.commentCache.secondary.SetAndUpdate(ctx, strconv.FormatInt(id, 10), res[])
		}()
		//
		return res, nil
	}
	return nil, err
}
