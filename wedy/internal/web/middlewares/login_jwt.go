package middlewares

import (
	"GoProj/wedy/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddleWareBuilder struct{}

func (m *LoginJWTMiddleWareBuilder) CheckJWTLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/user/login" || path == "/user/signup" || path == "/user/loginjwt" || path == "/user/loginsms/code/send" || path == "/user/loginsms/" {
			return
		}
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		signs := strings.Split(authCode, " ")
		if len(signs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := signs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.SigningKey, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		expireTime := uc.ExpiresAt
		if expireTime.Sub(time.Now()).Seconds() < 50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Second * 60))
			tokenStr, err = token.SignedString(web.SigningKey)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				log.Panicln(err)
			}
		}
		ctx.Set("user", uc)
	}
}
