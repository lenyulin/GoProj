package integration

import (
	service2 "GoProj/wedy/interactive/service"
	"GoProj/wedy/internal/repository"
	"GoProj/wedy/internal/repository/cache"
	dao "GoProj/wedy/internal/repository/dao/video"
	"GoProj/wedy/internal/service"
	"GoProj/wedy/internal/web"
	"GoProj/wedy/ioc"
	oss2 "GoProj/wedy/pkg/oss"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUploadNewVideoHandler(t *testing.T) {
	testcases := []struct {
		name       string
		before     func(t *testing.T)
		after      func(t *testing.T)
		art        Arcical
		wantCode   int
		wantResult Result[int64]
	}{
		{
			name: "new video",
			before: func(t *testing.T) {
				//准备数据
			},
			after: func(t *testing.T) {
				//验证已上传成功
			},
			art: Arcical{
				Title:   "My video title",
				Content: "My video content",
			},
			wantCode: http.StatusOK,
			wantResult: Result[int64]{
				Code: 1,
				Msg:  "Ok.",
			},
		},
	}

	oss := ioc.InitOSS()
	db := ioc.InitDB()
	redis := ioc.InitRedis()
	cache := cache.NewVideoRedisCache(redis)
	ossUploader := oss2.NewOSSHandler(oss)
	gormDAO := dao.NewGORMVideoDAO(db, cache)
	vdao := dao.NewOSSVideoDao(ossUploader, gormDAO)
	vrepo := repository.NewVideoRepository(vdao, cache)
	vSvc := service.NewVideoService(vrepo)
	interc := service2.NewInteractiveService()
	vhdl := web.NewVideoHandler(vSvc, interc)
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user", web.UserClaims{UserId: 1234567})
	})
	vhdl.RegisiterRoutes(server)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := os.Open("C:\\Users\\lyl69\\GolandProjects\\GoProj\\wedy\\internal\\integration\\media\\video.mp4")
	assert.NoError(t, err)
	mediaData, err := writer.CreateFormFile("media", "video.mp4")
	assert.NoError(t, err)
	_, err = io.Copy(mediaData, file)
	assert.NoError(t, err)
	reqJson, _ := json.Marshal(Req{
		Title:   "My video title",
		Content: "My video content",
	})
	err = writer.WriteField("json", string(reqJson))
	assert.NoError(t, err)
	err = writer.Close()
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			req, err := http.NewRequest(http.MethodPost, "/video/new", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.wantCode, recorder.Code)
			var res Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, res)
		})
	}
}
func TestVideoEditHandler(t *testing.T) {
	testcases := []struct {
		name       string
		before     func(t *testing.T)
		after      func(t *testing.T)
		art        Arcical
		wantCode   int
		wantResult Result[int64]
	}{
		{
			name: "edit video",
			before: func(t *testing.T) {
				//准备数据
			},
			after: func(t *testing.T) {
				//验证已上传成功
			},
			art: Arcical{
				Title:   "My new video title",
				Content: "My new video content",
			},
			wantCode: http.StatusOK,
			wantResult: Result[int64]{
				Code: 1,
				Msg:  "Edit Success",
			},
		},
	}

	oss := ioc.InitOSS()
	db := ioc.InitDB()
	gormDAO := dao.NewGORMVideoDAO(db)
	ossUploader := oss2.NewOSSHandler(oss)
	vdao := dao.NewOSSVideoDao(ossUploader, gormDAO)
	vrepo := repository.NewVideoRepository(vdao)
	vSvc := service.NewVideoService(vrepo)
	vhdl := web.NewVideoHandler(vSvc)
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user", web.UserClaims{UserId: 1234567})
	})
	vhdl.RegisiterRoutes(server)
	type Req struct {
		Id      int64  `json:"vuid"`
		Content string `json:"content"`
		Title   string `json:"title"`
	}
	var req Req
	req.Id = 3976732584221028914
	req.Content = "My new video content"
	req.Title = "My new video title"
	buff := &bytes.Buffer{}
	body, _ := json.Marshal(req)
	buff.Write(body)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			req, err := http.NewRequest(http.MethodPost, "/video/edit", buff)
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.wantCode, recorder.Code)
			var res Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, res)
		})
	}
}

type Arcical struct {
	Title   string
	Content string
}
type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
