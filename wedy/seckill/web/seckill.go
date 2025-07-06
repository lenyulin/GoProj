package web

import (
	"GoProj/wedy/pkg/bigcachex"
	"GoProj/wedy/pkg/bigcachex/proto"
	"GoProj/wedy/pkg/limiter"
	"GoProj/wedy/seckill/domain"
	"GoProj/wedy/seckill/service"
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
	svc     service.Seckill
}

func NewSeckillHandler(cache bigcachex.HybridCache, limiter limiter.Limiter, svc service.Seckill) *SeckillHandler {
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
		activityId    int64
		productId     int64
		quantity      int64
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
	acInfo, err := h.cache.Get(ctx, strconv.FormatInt(req.activityId, 10))
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
	if req.quantity <= 0 || req.quantity > activity.LimitPerUser {
		ctx.String(http.StatusOK, "Insufficient quantity")
		return
	}
	orderId, err := h.svc.Processing(ctx, domain.Order{
		UserId:        usr.UserId,
		ActivityId:    req.activityId,
		ProductId:     req.productId,
		PaymentMethod: req.paymentMethod,
		PromoCode:     req.promoCode,
	})
	if err != nil {
		ctx.String(http.StatusOK, "Internal Server Error")
		return
	}
	ctx.String(http.StatusOK, orderId)
	return
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
	res, err := h.svc.Cancel(ctx)
	if err != nil {
		ctx.String(http.StatusOK, "Internal Server Error")
		return
	}
	ctx.JSON(http.StatusOK, res)
	return
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
	if o.UserId != usr.UserId {
		ctx.String(http.StatusOK, "Insufficient user")
		return
	}
	orderStatus, err := h.svc.Status(ctx, o)
	if err != nil {
		ctx.String(http.StatusOK, "Internal Server Error")
		return
	}
	ctx.JSON(http.StatusOK, orderStatus)
	return
}
