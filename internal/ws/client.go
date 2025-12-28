package ws

import (
	"encoding/json"
	"ice-chat/internal/model/response"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn // WS连接
	send     chan []byte     // 消息发送通道
	username string          // 用户名
}

func NewClient(conn *websocket.Conn, username string, sendBufferSize int) *Client {
	return &Client{
		conn:     conn,
		send:     make(chan []byte, sendBufferSize),
		username: username,
	}
}

func (c *Client) ReadPump(manager *ClientManager, msgHandler func([]byte, *Client)) {
	defer func() {
		// 客户端发生异常，直接断开
		manager.unregister <- c
		_ = c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取消息错误：%v", err)
			}
			break
		}

		// 调用外部消息处理器
		msgHandler(msg, c)
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// 通道已经关闭了，直接中断当前的 websocket 连接
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 这里源码写法
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err = w.Write(message); err != nil {
				return
			}
			if err := w.Close(); err != nil {
				return
			}
		default:
			// 心跳检测机制
		}
	}
}

func (c *Client) SendErrMessageToClient(errMsg string) error {
	var errorMessage = response.Message{
		Error:   true,
		Content: errMsg,
	}

	emb, err := json.Marshal(errorMessage)

	if err != nil {
		return err
	}

	c.send <- emb
	return nil
}

func (c *Client) Close() {
	_ = c.conn.Close()
	close(c.send)
}
