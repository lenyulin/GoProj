package web

import (
	service2 "GoProj/wedy/interactive/service"
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/service"
	"GoProj/wedy/pkg/snowflake"
	"encoding/json"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

const Dst = "C:\\Users\\lyl69\\GolandProjects\\GoProj\\wedy\\tmp\\"

type VideoHandler struct {
	svc      service.VideoService
	interSvc service2.InteractiveService
	biz      string
}

func NewVideoHandler(svc service.VideoService, interSvc service2.InteractiveService) *VideoHandler {
	return &VideoHandler{
		svc:      svc,
		interSvc: interSvc,
		biz:      "video",
	}
}
func (v *VideoHandler) RegisiterRoutes(server *gin.Engine) {
	vg := server.Group("/video")
	vg.POST("/edit", v.Edit)
	vg.POST("/new", v.Publish)
	vg.POST("/withdrawn", v.Withdrawn)
	vg.GET("/detail/:id", v.Detail)
	vg.POST("/list", v.List)
}
func (v *VideoHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		VUid    int64  `json:"vuid"`
		Status  uint8  `json:"status"`
	}
	var req Req
	uc := ctx.MustGet("user").(UserClaims)
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Results[int16]{
			Code: 1,
			Msg:  "Internal Error",
		})
		return
	}
	err := v.svc.Update(ctx, domian.Video{
		Title:   req.Title,
		Content: req.Content,
		VUid:    req.VUid,
		Uid:     uc.UserId,
		Status:  domian.VideoStatus(req.Status),
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Results[int16]{
			Code: 1,
			Msg:  "Internal Error",
		})
	}
	ctx.JSON(http.StatusOK, Results[int16]{
		Code: 1,
		Msg:  "Edit Success",
	})
}
func (v *VideoHandler) Publish(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	jsonReq := ctx.Request.FormValue("json")
	var req Req
	err := json.Unmarshal([]byte(jsonReq), &req)
	if err != nil {
		ctx.String(http.StatusOK, "Internal Error")
		return
	}
	uc := ctx.MustGet("user").(UserClaims)
	//uc := int64(123123123)
	media, err := ctx.FormFile("media")
	if err != nil {
		ctx.String(http.StatusOK, "Internal Error")
		return
	}
	if media.Size > 64<<20 {
		ctx.String(http.StatusOK, "Media too large")
		return
	}
	uid, err := snowflake.Generate()
	if err != nil {
		ctx.String(http.StatusOK, "Internal Error")
		return
	}
	fileDst := fmt.Sprintf("%s%d.mp4", Dst, uid)
	if err := ctx.SaveUploadedFile(media, fileDst); err != nil {
		ctx.String(http.StatusOK, "Interal Error")
		return
	}

	err = v.svc.Publish(ctx, domian.Video{
		Title:   req.Title,
		Content: req.Content,
		VUid:    uid,
		Uid:     uc.UserId,
		Author: domian.Author{
			Id: uc.UserId,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Results[int64]{
			Code: 2,
			Msg:  "Interal Error",
		})
	}
	if err := os.Remove(fileDst); err != nil {
		println(err)
	}
	ctx.JSON(http.StatusOK, Results[int64]{
		Code: 1,
		Msg:  "Ok.",
	})
}

func (v *VideoHandler) Withdrawn(ctx *gin.Context) {
	type Req struct {
		VUid int64 `json:"vuid"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "Interal Error")
	}
	uc := ctx.MustGet("user").(UserClaims)
	err := v.svc.Withdrawn(ctx, domian.Video{
		VUid: req.VUid,
		Uid:  uc.UserId,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Results[int64]{
			Code: 2,
			Msg:  "Interal Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, Results[int64]{
		Code: 1,
		Msg:  "Withdrawn Success",
	})
}

func (v *VideoHandler) Detail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Results[int64]{
			Code: 1,
			Msg:  "Id Error",
		})
		return
	}
	uc := ctx.MustGet("user").(UserClaims)
	videos, err := v.svc.GetById(ctx, id, uc.UserId)
	if err != nil {
		ctx.JSON(http.StatusOK, Results[int64]{
			Code: 1,
			Msg:  "Internal Error",
		})
		return
	}
	if uc.UserId == videos.Uid {
		ctx.JSON(http.StatusOK, Results[int64]{
			Code: 1,
			Msg:  "Internal Error",
		})
	}
	vo := VideoVO{
		Title:    videos.Title,
		Content:  videos.Content,
		Id:       videos.VUid,
		AuthorId: videos.Author.Id,
	}
	ctx.JSON(http.StatusOK, Results[int64]{
		Data: vo,
	})
}
func (v *VideoHandler) PubDetail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Results[int64]{
			Code: 1,
			Msg:  "Id Error",
		})
		return
	}
	uc := ctx.MustGet("user").(UserClaims)
	videos, err := v.svc.GetById(ctx, id, uc.UserId)
	if err != nil {
		ctx.JSON(http.StatusOK, Results[int64]{
			Code: 1,
			Msg:  "Internal Error",
		})
		return
	}
	err = v.interSvc.IncrReadCnt(ctx, v.biz, id)
	ctx.JSON(http.StatusOK, Results[int64]{
		Data: videos,
	})
}
func (v *VideoHandler) List(ctx *gin.Context) {
	var page Page
	if err := ctx.Bind(&page); err != nil {
		ctx.JSON(http.StatusOK, Results[int64]{
			Code: 1,
			Msg:  "Internal Error",
		})
		return
	}
	uc := ctx.MustGet("user").(UserClaims)
	videos, err := v.svc.GetByAuthor(ctx, uc.UserId, page.Limit, page.Offset)
	if err != nil {
		ctx.JSON(http.StatusOK, Results[int64]{
			Code: 1,
			Msg:  "Internal Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, Results[int64]{
		Data: slice.Map[domian.Video, VideoVO](videos, func(idx int, src domian.Video) VideoVO {
			return VideoVO{
				Title:    src.Title,
				Content:  src.Content,
				Id:       src.VUid,
				AuthorId: src.Author.Id,
			}
		}),
	})
}

type Results[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
