package web

import (
	"GoProj/wedy/comment/domain"
	"GoProj/wedy/comment/service"
	svcmock "GoProj/wedy/internal/service/mocks"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommentHandler_Submit(t *testing.T) {
	testcases := []struct {
		name     string
		wandCode int
		reqBody  string
		wantResp CommentResult[int64]
		mock     func(ctrl *gomock.Controller) service.CommentService
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) service.CommentService {
				svc := svcmock.NewMockCommentService(ctrl)
				//svc := svcmock.NewMockCommentService(ctrl)
				svc.EXPECT().SubmitComment(gomock.Any(), domain.Comment{
					User:    1234567,
					Uid:     1234567,
					VUid:    2345678,
					Content: "hello",
				}).Return(nil)
				return svc
			},
			reqBody:  `{"user":1234567,"uid":1234567,"vuid":2345678,"content":"hello"}`,
			wandCode: http.StatusOK,
			wantResp: CommentResult[int64]{
				Code: 1,
				Msg:  "Ok.",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commentService := tc.mock(ctrl)
			hel := NewCommentHandler(commentService)
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user", UserClaims{UserId: 1234567})
			})
			hel.RegisiterRoutes(server)
			req, _ := http.NewRequest(http.MethodPost, "/comment/submit", bytes.NewBufferString(tc.reqBody))
			req.Header.Set("content-type", "application/json")
			recoder := httptest.NewRecorder()
			server.ServeHTTP(recoder, req)

			var res CommentResult[int64]
			err := json.NewDecoder(recoder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wandCode, recoder.Code)
			assert.Equal(t, tc.wantResp, res)

		})
	}
}

type CommentResult[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:",data"`
}
