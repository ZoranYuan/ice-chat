package api

import (
	"ice-chat/internal/response"

	"github.com/gin-gonic/gin"
)

func HeartCheckApi(ctx *gin.Context) {
	response.OKWithMsg(ctx, "pong")
}
