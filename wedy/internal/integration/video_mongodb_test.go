package integration

import (
	"GoProj/wedy/internal/integration/startup"
	"GoProj/wedy/internal/repository"
	"GoProj/wedy/internal/repository/dao/video"
	"GoProj/wedy/internal/service"
	"GoProj/wedy/internal/web"
	"GoProj/wedy/ioc"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type VideoMongoDBHandelSuite struct {
	suite.Suite
	server  *gin.Engine
	mdb     *mongo.Database
	col     *mongo.Collection
	bucket  *gridfs.Bucket
	gormDAO *dao.GORMVideoDAO
}

func (s *VideoMongoDBHandelSuite) SetupSuite() {
	mdb := startup.InitMongoDB()
	s.mdb = mdb
	s.col = mdb.Collection("video")
	bucket, err := gridfs.NewBucket(s.mdb)
	assert.NoError(s.T(), err)
	s.bucket = bucket
	db := ioc.InitDB()
	s.gormDAO = dao.NewGORMVideoDAO(db)
	vdao := dao.NewMongoDBDAO(s.mdb, s.col, s.bucket, s.gormDAO)
	vrepo := repository.NewVideoRepository(vdao)
	vSvc := service.NewVideoService(vrepo)
	vhdl := web.NewVideoHandler(vSvc)
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user", web.UserClaims{UserId: 1234567})
	})
	vhdl.RegisiterRoutes(server)
	s.server = server
}

type Req struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (s *VideoMongoDBHandelSuite) TestVideoPublish() {
	t := s.T()
	testcases := []struct {
		name       string
		before     func(t *testing.T)
		after      func(t *testing.T)
		art        Req
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
				filter := bson.D{{"metadata.AuthorId", 1234567}}
				cursor, err := s.bucket.Find(filter)
				assert.NoError(t, err)
				defer func() {
					if err := cursor.Close(context.TODO()); err != nil {
						log.Panic(err)
					}
				}()
				type gridfsFile struct {
					Name   string             `bson:"filename"`
					Length int64              `bson:"length"`
					Meta   dao.Video          `bson:"metadata"`
					id     primitive.ObjectID `bson:"$oid"`
				}
				var foundFiles []gridfsFile
				err = cursor.All(context.TODO(), &foundFiles)
				assert.NoError(t, err)
				for _, res := range foundFiles {
					assert.True(t, res.Meta.Ctime > 0)
					assert.True(t, res.Meta.Utime > 0)
					assert.Equal(t, res.Meta.Content, "My video content")
					assert.Equal(t, res.Meta.Title, "My video title")
					err = s.bucket.Delete(foundFiles[0].id)
					assert.NoError(t, err)
				}
			},
			art: Req{
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
			s.server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.wantCode, recorder.Code)
			var res Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, res)
		})
	}
}
func TestArticleMongoDBHandler(t *testing.T) {
	suite.Run(t, &VideoMongoDBHandelSuite{})
}
