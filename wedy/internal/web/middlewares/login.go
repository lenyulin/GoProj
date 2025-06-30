package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddleWareBuilder struct{}

func (m *LoginMiddleWareBuilder) CheckLoginSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/user/login" || path == "/user/signup" || path == "/user/loginjwt" {
			return
		}
		session := sessions.Default(ctx)
		if session.Get("userId") == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		return
	}
}
