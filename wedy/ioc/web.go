package ioc

import (
	"GoProj/wedy/internal/web"
	"GoProj/wedy/internal/web/middlewares"
	"GoProj/wedy/pkg/ginx"
	"GoProj/wedy/pkg/ginx/middlewares/prometheus"
	ratelimit "GoProj/wedy/pkg/ginx/middlewares/ratelimite"
	"GoProj/wedy/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UeerHandler, videoHdl *web.VideoHandler, comHdl *web.CommentHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisiterRoutes(server)
	videoHdl.RegisiterRoutes(server)
	comHdl.RegisiterRoutes(server)
	return server
}

func InitMiddleware(redisClient redis.Cmdable) []gin.HandlerFunc {
	pb := &prometheus.Builder{
		NameSpace: "lenyulin",
		Subsystem: "wedy",
		Name:      "gin_http",
	}
	ginx.InitCounter(prometheus2.CounterOpts{
		Namespace: "lenyulin",
		Subsystem: "wedy",
		Name:      "code",
	})
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return strings.Contains(origin, "yulin_lei.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		otelgin.Middleware("wedy"),
		pb.BuildResponseTime(),
		pb.BuildActiveRequest(),
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 10)).Build(),
		(&middlewares.LoginJWTMiddleWareBuilder{}).CheckJWTLogin(),
	}
}
