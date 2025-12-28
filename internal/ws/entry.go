package ws

import (
	"ice-chat/config"
	"net/http"

	"github.com/gorilla/websocket"
)

func GetUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  config.Conf.Ws.ReadBufferSize,
		WriteBufferSize: config.Conf.Ws.ReadBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}
