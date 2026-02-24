package router

import (
	"ice-chat/internal/api"
	"ice-chat/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoomsRouter(r *gin.Engine, roomsApi api.RoomsApi) {
	gr := r.Group("/room").Use(middleware.AuthMiddleware())
	gr.POST("/create", roomsApi.CreateRoom)
	gr.POST("/join/:joinCode", roomsApi.JoinRoom)
	gr.GET("/watch/:roomId", roomsApi.GetWatchVideo)
}

// TODO 当前端建立一起看的 ws 服务时
/*
首先：同步当前群聊一起看的状态
然后，各个用户自行播放即可
当用户拖动进度条时，通过 ws 广播给其他用户即可
*/
