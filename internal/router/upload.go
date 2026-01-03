package router

import (
	"ice-chat/internal/api"
	"ice-chat/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUploadRouter(r *gin.Engine, uploadApi api.UploadApi) {
	ur := r.Group("/upload")
	ur.Use(middleware.AuthMiddleware())

	ur.POST("", uploadApi.Upload)
	ur.POST("/merge", uploadApi.Merge)
	ur.POST("/init", uploadApi.UploadInit)
}
