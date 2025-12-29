package router

import (
	"ice-chat/internal/api"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine, userApi *api.UserAPI) {
	ur := r.Group("/user")
	ur.POST("/login", userApi.Login)
	ur.POST("/upload", userApi.Upload)
}
