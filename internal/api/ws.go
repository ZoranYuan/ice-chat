package api

import (
	"ice-chat/config"
	"ice-chat/internal/response"
	"ice-chat/internal/service"
	"ice-chat/internal/ws"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsAPI struct {
	wsSvc     *service.WsService
	wsManager *ws.ClientManager
	wsUtils   *websocket.Upgrader
}

func NewWsAPI(wsSvc *service.WsService, wsManager *ws.ClientManager, wsUtils *websocket.Upgrader) *WsAPI {
	return &WsAPI{
		wsSvc:     wsSvc,
		wsManager: wsManager,
		wsUtils:   wsUtils,
	}
}

func (w *WsAPI) Chat(ctx *gin.Context) {
	userName := ctx.Query("userName")
	if userName == "" {
		response.BadRequest(ctx)
		ctx.Abort()
		return
	}

	conn, err := w.wsUtils.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WS升级失败：%v", err)
		response.InternalError(ctx, "failed to establish websocket service")
		return
	}

	client := ws.NewClient(conn, userName, config.Conf.Ws.ReadBufferSize)
	w.wsManager.AddClient(client)
	go client.WritePump()
	go client.ReadPump(w.wsManager, w.wsSvc.MessageHandler)
}
