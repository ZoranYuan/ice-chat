package ws

import (
	"fmt"
	"sync"
)

type Room struct {
	id         uint64
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte

	mu sync.RWMutex
}

func NewRoom(roomId uint64) *Room {
	return &Room{
		id:         roomId,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 10),
		broadcast:  make(chan []byte, 256),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.mu.Lock()
			r.clients[client] = true
			r.mu.Unlock()
		case client := <-r.unregister:
			r.mu.Lock()
			delete(r.clients, client)
			close(client.send)
			r.mu.Unlock()
			if len(r.clients) == 0 {
				// TODO 交给 Hub 决定是否删除
			}
		case message := <-r.broadcast:
			for client := range r.clients {
				fmt.Print("start to boradcast message")
				select {
				case client.send <- message:
				default:
					// 客户端阻塞，直接踢掉
					close(client.send)
					client.conn.Close()
					delete(r.clients, client)
				}
			}
		}
	}
}

func (r *Room) AddClient(c *Client) {
	r.register <- c
}

func (r *Room) RemoveClient(c *Client) {
	r.unregister <- c
	c.conn.Close()
	close(c.send)
}

func (r *Room) Broadcast(msg []byte) {
	r.broadcast <- msg
}
