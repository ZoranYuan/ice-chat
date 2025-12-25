package router

import (
	"ice-chat/internal/api"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine, userApi *api.UserAPI) {
	r.POST("/user/login", userApi.Login)
}
