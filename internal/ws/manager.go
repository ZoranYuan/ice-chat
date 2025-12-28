package ws

import (
	"sync"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client // 注销通道
	mu         sync.RWMutex // 读写锁
}

func NewManager() *ClientManager {
	return &ClientManager{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// 开启一个协程进行 ClientManager 的 Run
func (cm *ClientManager) Run() {
	for {
		select {
		case client := <-cm.register:
			cm.mu.Lock()
			cm.clients[client] = true
			cm.mu.Unlock()
		case client := <-cm.unregister:
			cm.mu.Lock()
			if _, ok := cm.clients[client]; ok {
				delete(cm.clients, client)
				close(client.send)
			}
			cm.mu.Unlock()
		case message := <-cm.broadcast:
			for client := range cm.clients {
				// 对于某个 client 而言，写入可能会造成阻塞，第一是缓冲区的问题，第二就是前端连接 done 了
				// TODO 后续处理
				client.send <- message
			}
		}
	}
}

// Broadcast 广播消息（外部调用）
func (cm *ClientManager) Broadcast(msg []byte) {
	cm.broadcast <- msg
}

// GetOnlineCount 获取在线人数
func (m *ClientManager) GetOnlineCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}

func (m *ClientManager) AddClient(client *Client) {
	m.register <- client
}
