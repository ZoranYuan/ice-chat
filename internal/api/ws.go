package api

import (
	"ice-chat/config"
	"ice-chat/internal/response"
	"ice-chat/internal/service"
	"ice-chat/pkg/ws"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type wsApiImpl struct {
	wsSvc       service.WsService
	roomManager *ws.RoomManager
	wsUtils     *ws.WsUtils
}

type WsApi interface {
	Chat(ctx *gin.Context)
	Watch(ctx *gin.Context)
}

func NewWsAPI(wsSvc service.WsService, roomManager *ws.RoomManager, wsUtils *ws.WsUtils) WsApi {
	return &wsApiImpl{
		wsSvc:       wsSvc,
		roomManager: roomManager,
		wsUtils:     wsUtils,
	}
}

func (w *wsApiImpl) handle(ctx *gin.Context, handleFn func(c *ws.Client, room *ws.Room, msg []byte)) {
	v := ctx.Param("groupId")
	uid, exists := ctx.Get("uid")
	if !exists {
		response.Unauthorized(ctx)
		ctx.Abort()
		return
	}

	groupId, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		response.BadRequestWithMessage(ctx, "参数错误")
		ctx.Abort()
		return
	}

	if exists, _ := w.wsSvc.GroupIsExists(groupId); !exists {
		response.BadRequestWithMessage(ctx, "房间号不存在")
		// TODO 可以额外处理下 err
		ctx.Abort()
		return
	}

	conn, err := w.wsUtils.GetUpgrader().Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WS升级失败：%v", err)
		response.InternalError(ctx, "failed to establish websocket service")
		return
	}

	// 获取 room
	client := ws.NewClient(conn, uid.(uint64), config.Conf.Ws.WriteBufferSize)
	room := w.roomManager.GetRoom(groupId)
	room.AddClient(client)
	go client.Write(room)
	go client.Read(room, handleFn)
}

func (w *wsApiImpl) Chat(ctx *gin.Context) {
	w.handle(ctx, w.wsSvc.ChatHandler)
}

func (w *wsApiImpl) Watch(ctx *gin.Context) {
	w.handle(ctx, w.wsSvc.WatchHandler)
}
