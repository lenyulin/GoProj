package web

import (
	"GoProj/wedy/pkg/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
)

type SeckillHandler struct {
	redis   redis.Cmdable
	limiter limiter.Limiter
}

func NewSeckillHandler(redis redis.Cmdable) *SeckillHandler {
	return &SeckillHandler{
		redis: redis,
	}
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId    int64
	UserAgent string
}

func (h *SeckillHandler) RegisiterRoutes(server *gin.Engine) {
	server.POST("/seckill", h.SecKill)
}
func (h *SeckillHandler) SecKill(ctx *gin.Context) {
	type ac struct {
		activity_id string
		product_id  string
		promo_code  string
	}
	usr, ok := ctx.MustGet("user").(UserClaims)
	if !ok {
		ctx.String(http.StatusBadRequest, "Internal Server Error")
		return
	}
	if limited, _ := h.limiter.Limit(ctx, strconv.FormatInt(usr.UserId, 10)); limited {
		ctx.String(http.StatusBadRequest, "Too many requests")
		return
	}
	var req ac
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusBadRequest, "Internal Server Error")
		return
	}

}
