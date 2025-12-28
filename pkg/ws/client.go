package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
	uid  uint64
}

func NewClient(conn *websocket.Conn, uid uint64, sendBufferSize int) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, sendBufferSize),
		uid:  uid,
	}
}

func (c *Client) Read(room *Room, handle func(c *Client, room *Room, msg []byte)) {
	for {
		defer func() {
			if r := recover(); r != nil {
				room.RemoveClient(c)
				c.conn.Close()
				close(c.send)
			}
		}()

		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取消息错误：%v", err)
				panic("close")
			}
			break
		}

		// TODO 将消息进行进一步过滤
		handle(c, room, msg)
	}
}

func (c *Client) Write(room *Room) {
	defer func() {
		close(c.send)
		c.conn.Close()
		room.RemoveClient(c)
	}()

	for {
		msg, ok := <-c.send
		if !ok {
			// 通道已经关闭了，直接中断当前的 websocket 连接
			_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		if _, err = w.Write(msg); err != nil {
			return
		}
		if err := w.Close(); err != nil {
			return
		}
	}
}
