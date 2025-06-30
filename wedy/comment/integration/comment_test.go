package integration

import (
	"GoProj/wedy/comment/repository"
	"GoProj/wedy/comment/repository/dao"
	"GoProj/wedy/comment/service"
	"GoProj/wedy/internal/web"
	"GoProj/wedy/ioc"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestComment(t *testing.T) {
	testcases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		wantCode   int
		wantResult CommentResult[int64]

		reqBuilder func(t *testing.T) *http.Request
	}{
		{
			name: "success",
			reqBuilder: func(t *testing.T) *http.Request {
				req, _ := http.NewRequest(
					http.MethodPost,
					"/comment/submit",
					bytes.NewBuffer([]byte("{\"uid\": \"1234567\",\"vuid\": \"2345678\",\"content\": \"Content\"}")))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			before:   func(t *testing.T) {},
			after:    func(t *testing.T) {},
			wantCode: http.StatusOK,
			wantResult: CommentResult[int64]{
				Code: 1,
				Msg:  "Ok.",
			},
		},
	}
	db := ioc.InitDB()
	cdao := dao.NewGORMCommentDAO(db)
	crepo := repository.NewCachedCommentRepository(cdao)
	cSvc := service.NewCommentService(crepo)
	chdl := web.NewCommentHandler(cSvc)
	server := gin.Default()
	chdl.RegisiterRoutes(server)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			recoder := httptest.NewRecorder()
			server.ServeHTTP(recoder, tc.reqBuilder(t))
			var res CommentResult[int64]
			err := json.Unmarshal(recoder.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, recoder.Code)
			assert.Equal(t, tc.wantResult, res)
		})
	}
}

type CommentResult[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:",data"`
}
