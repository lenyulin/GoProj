package web

import (
	"GoProj/wedy/pkg/bigcachex"
	"GoProj/wedy/pkg/bigcachex/proto"
	"GoProj/wedy/pkg/limiter"
	"GoProj/wedy/seckill-system/domain"
	"GoProj/wedy/seckill-system/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	proto2 "google.golang.org/protobuf/proto"
	"net/http"
	"strconv"
	"time"
)

type SeckillHandler struct {
	cache   bigcachex.HybridCache
	limiter limiter.Limiter
	svc     service.OrderService
}

func NewSeckillHandler(cache bigcachex.HybridCache, limiter limiter.Limiter, svc service.OrderService) *SeckillHandler {
	return &SeckillHandler{
		cache:   cache,
		limiter: limiter,
		svc:     svc,
	}
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId    int64
	UserAgent string
}

func (h *SeckillHandler) RegisiterRoutes(server *gin.Engine) {
	sg := server.Group("/seckill")
	sg.POST("/order-placement", h.Placement)
	sg.POST("/order-cancel", h.Cancel)
	sg.POST("/order-status", h.Status)
}

func (h *SeckillHandler) Placement(ctx *gin.Context) {
	type ac struct {
		activityId    string
		productId     string
		paymentMethod string
		promoCode     []string
	}
	usr, ok := ctx.MustGet("user").(UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "Unauthorized User")
		return
	}
	if limited, _ := h.limiter.Limit(ctx, strconv.FormatInt(usr.UserId, 10)); limited {
		ctx.String(http.StatusOK, "Too many requests")
		return
	}
	var req ac
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "Internal Server Error")
		return
	}
	acInfo, err := h.cache.Get(ctx, req.activityId)
	if err != nil {
		ctx.String(http.StatusOK, "Internal Server Error")
		return
	}
	var activity proto.SeckillActivity
	err = proto2.Unmarshal(acInfo, &activity)
	if err != nil {
		ctx.String(http.StatusOK, "Internal Server Error")
		return
	}
	if activity.StartTime > time.Now().UnixMilli() || activity.EndTime < time.Now().UnixMilli() || activity.Status != 2 {
		ctx.String(http.StatusOK, "Activity status: not started or ended")
		return
	}
	err = h.svc.Processing(ctx)
	if err != nil {
		ctx.String(http.StatusOK, "Internal Server Error")
		return
	}
}

// 取消订单

func (h *SeckillHandler) Cancel(ctx *gin.Context) {
	usr, ok := ctx.MustGet("user").(UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "Unauthorized User")
		return
	}
	if limited, _ := h.limiter.Limit(ctx, strconv.FormatInt(usr.UserId, 10)); limited {
		ctx.String(http.StatusOK, "Too many requests")
		return
	}
	err := h.svc.Cancel(ctx)
	if err != nil {
		ctx.String(http.StatusOK, "Internal Server Error")
		return
	}
}

//查询订单状态

func (h *SeckillHandler) Status(ctx *gin.Context) {
	usr, ok := ctx.MustGet("user").(UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "Unauthorized User")
		return
	}
	if limited, _ := h.limiter.Limit(ctx, strconv.FormatInt(usr.UserId, 10)); limited {
		ctx.String(http.StatusOK, "Too many requests")
		return
	}
	var o domain.Order
	err := ctx.Bind(&o)
	if err != nil {
		ctx.String(http.StatusOK, "Internal Server Error")
		return
	}
	o.UserId = usr.UserId
	err = h.svc.Status(ctx, o)
	if err != nil {
		ctx.String(http.StatusOK, "Internal Server Error")
		return
	}
}
