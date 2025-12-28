package router

import (
	"ice-chat/internal/api"
	"ice-chat/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterWsRouter(r *gin.Engine, wsApi api.WsApi) {
	wr := r.Group("/ws").Use(middleware.AuthMiddleware())
	wr.GET("/chat/:groupId", wsApi.Chat)
	wr.GET("/watch/:groupId", wsApi.Watch)
}
