package router

import (
	"ice-chat/internal/api"
	"ice-chat/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterGroupsRouter(r *gin.Engine, groupsApi api.GroupApi) {
	gr := r.Group("/group").Use(middleware.AuthMiddleware())
	gr.POST("/create", groupsApi.Create)
}
