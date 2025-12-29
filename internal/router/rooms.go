package router

import (
	"ice-chat/internal/api"
	"ice-chat/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoomsRouter(r *gin.Engine, roomsApi api.RoomsApi) {
	gr := r.Group("/room").Use(middleware.AuthMiddleware())
	gr.POST("/create", roomsApi.Create)
}
