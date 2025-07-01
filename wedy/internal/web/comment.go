package web

import (
	"GoProj/wedy/comment/domain"
	"GoProj/wedy/comment/service"
	"GoProj/wedy/internal/domian"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ErrBindCommentError = errors.New("Interal Error")

//	type Comment struct {
//		Uid     string `json:"uid"`
//		VUid    string `json:"vuid"`
//		Content string `json:"content"`
//	}
type CommentHandler struct {
	svc service.CommentService
}

func (h *CommentHandler) RegisiterRoutes(server *gin.Engine) {
	ug := server.Group("/comment")
	ug.POST("/submit", h.Submit)
}
func NewCommentHandler(svc service.CommentService) *CommentHandler {
	return &CommentHandler{
		svc: svc,
	}
}

type Comment struct {
	User    int64  `json:"user"`
	Id      int64  `json:"id"`
	Content string `json:"content"`
	Vid     int64  `json:"vid"`
}

func (h *CommentHandler) Submit(ctx *gin.Context) {
	var comment Comment
	err := ctx.Bind(&comment)
	if err != nil {
		ctx.String(http.StatusOK, "Internal Error")
		return
	}
	usr := ctx.MustGet("user").(UserClaims)
	err = h.svc.SubmitComment(ctx, domain.Comment{
		User:    &domian.User{Id: usr.UserId},
		Id:      comment.Id,
		Content: comment.Content,
		Vid:     comment.Vid,
	})
	if err != nil {
		ctx.String(http.StatusOK, ErrBindCommentError.Error())
		return
	}
	ctx.JSON(
		http.StatusOK,
		Results[int64]{
			Code: 1,
			Msg:  "Ok.",
		},
	)
}
