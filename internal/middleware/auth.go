package middleware

import (
	"ice-chat/config"
	"ice-chat/internal/response"
	"ice-chat/pkg/jwtUtils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		if token == "" {
			response.Unauthorized(ctx)
			ctx.Abort()
			return
		}

		ju := jwtUtils.CreateJwtUtils(jwtUtils.Config{
			Secret: []byte(config.Conf.JWT.AccessTokenExpireDuration),
			Expire: config.Conf.JWT.GetAccessTokenExpireDuration(),
		})
		token = strings.TrimPrefix(token, "Bearer ")
		// TODO 根据业务需求去截取字段
		claims, err := ju.Parse(token)
		if err != nil {
			response.Unauthorized(ctx)
			ctx.Abort()
			return
		}

		// 注入上下文
		ctx.Set("uid", claims.UserId)
		ctx.Set("jti", claims.JTI)

		ctx.Next()
	}
}
