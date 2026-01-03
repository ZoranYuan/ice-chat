package ws

import (
	"fmt"
	"log"
	"runtime/debug"

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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Read goroutine panic: %v\n%s", r, debug.Stack())
		}
	}()

	for {
		var msg []byte
		var err error

		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("ReadMessage panic: %v\n%s", r, debug.Stack())
					err = fmt.Errorf("panic during ReadMessage")
				}
			}()

			_, msg, err = c.conn.ReadMessage()
		}()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取消息错误: %v", err)
			}
			break // 出现错误或连接关闭，退出循环
		}

		if room != nil && msg != nil {
			safeHandle(c, room, msg, handle)
		}
	}

	log.Println("Read goroutine exited")
}

// 安全调用 handle 的包装函数，捕获 handle 内部 panic
func safeHandle(c *Client, room *Room, msg []byte, handle func(c *Client, room *Room, msg []byte)) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("handle panic: %v\n%s", r, debug.Stack())
		}
	}()

	handle(c, room, msg)
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
