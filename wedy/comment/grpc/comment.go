package grpc

import (
	commv1 "GoProj/wedy/api/proto/gen/wedy/api/proto/comm/v1"
	"GoProj/wedy/comment/domain"
	"GoProj/wedy/comment/service"
	"context"
	"errors"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

var (
	ErrSubmitCommentFailed = errors.New("submit comment failed")
	ErrGetCommentFailed    = errors.New("get comment failed")
	ErrLikeCommentFailed   = errors.New("increase like comment failed")
)

type CommentServiceServer struct {
	commv1.UnimplementedCommentServiceServer
	svc service.CommentService
}

func NewCommentServiceServer(svc service.CommentService) *CommentServiceServer {
	return &CommentServiceServer{svc: svc}
}

func (s *CommentServiceServer) Register(server *grpc.Server) {
	commv1.RegisterCommentServiceServer(server, s)
}
func (s *CommentServiceServer) Like(ctx context.Context, req *commv1.LikeRequest) (*commv1.LikeResponse, error) {
	err := s.svc.Like(ctx, req.Id, req.VId, req.Uid)
	if err != nil {
		return &commv1.LikeResponse{}, ErrLikeCommentFailed
	}
	return &commv1.LikeResponse{}, nil
}
func (s *CommentServiceServer) GetComment(ctx context.Context, req *commv1.GetCommentRequest) (*commv1.GetCommentResponse, error) {
	res, err := s.svc.Get(ctx, req.Vid, req.Page)
	if err != nil {
		return &commv1.GetCommentResponse{}, ErrGetCommentFailed
	}
	comm := slice.Map[domain.Comment, *commv1.Comments](res, func(idx int, src domain.Comment) *commv1.Comments {
		mus := slice.Map[domain.MentionedUser, *commv1.MentionedUser](src.MentionedUsers, func(idx int, sr domain.MentionedUser) *commv1.MentionedUser {
			return &commv1.MentionedUser{
				UserId:   sr.Id,
				UserName: sr.Name,
			}
		})
		return &commv1.Comments{
			Id:         src.Id,
			VId:        src.VId,
			Content:    src.Content,
			Ctime:      src.Ctime,
			Like:       src.Like,
			Dislike:    src.Dislike,
			ReplyCount: src.ReplyCount,
			Picture:    src.Picture,
			User: &commv1.User{
				Id:   src.User.Id,
				Name: src.User.Name,
			},
			MentionedUsers: mus,
		}
	})
	return &commv1.GetCommentResponse{Comm: comm}, err
}
func (s *CommentServiceServer) SubmitComment(ctx context.Context, req *commv1.SubmitCommentRequest) (*commv1.SubmitCommentResponse, error) {
	mus := slice.Map[*commv1.MentionedUser, domain.MentionedUser](req.MentionedUsers, func(idx int, src *commv1.MentionedUser) domain.MentionedUser {
		return domain.MentionedUser{
			Id:   src.UserId,
			Name: src.UserName,
		}
	})
	err := s.svc.Submit(ctx, domain.Comment{
		User: domain.User{
			Id:   req.UserId,
			Name: req.UserName,
		},
		VId:            req.Vid,
		Content:        req.Content,
		Picture:        req.Picture,
		MentionedUsers: mus,
	})
	if err != nil {
		return &commv1.SubmitCommentResponse{}, ErrSubmitCommentFailed
	}
	return &commv1.SubmitCommentResponse{}, nil
}
