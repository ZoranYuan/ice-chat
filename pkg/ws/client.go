package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	uid      uint64
	isClosed bool
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
				log.Printf("发生错误 %v", r)
			}
		}()

		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取消息错误：%v", err)
			}
			break
		}

		// TODO 将消息进行进一步过滤
		handle(c, room, msg)
	}
}

func (c *Client) Write(room *Room) {
	defer func() {
		_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		room.RemoveClient(c)
	}()

	for {
		msg, ok := <-c.send

		if !ok {
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

func (c *Client) GetUid() uint64 {
	return c.uid
}

func (c *Client) SendMessageToClient(msg []byte) {
	c.send <- msg
}
