package router

import (
	"ice-chat/internal/api"

	"github.com/gin-gonic/gin"
)

func RegisterHeartCheckRouter(r *gin.Engine) {
	r.GET("/ping", api.HeartCheckApi)
}
