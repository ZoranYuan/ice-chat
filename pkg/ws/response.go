package ws

import (
	"encoding/json"
	"log"
)

type WsResponse struct {
	Success bool
	Data    any
	Msg     string
}

func Ok(c *Client, data any) {
	var WsResponse = WsResponse{
		Success: true,
		Data:    data,
		Msg:     "数据同步成功",
	}

	WsResponseBytes, err := json.Marshal(WsResponse)
	if err != nil {
		log.Printf("failed to load video state: %v", err)
		return
	}
	c.SendMessageToClient(WsResponseBytes)
}

func Fail(c *Client, msg string) {
	var WsResponse = WsResponse{
		Success: false,
		Msg:     msg,
		Data:    nil,
	}

	WsResponseBytes, err := json.Marshal(WsResponse)
	if err != nil {
		log.Printf("failed to load video state: %v", err)
		return
	}
	c.SendMessageToClient(WsResponseBytes)
}
