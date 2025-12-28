package ws

import "sync"

type RoomManager struct {
	rooms map[uint64]*Room
	mu    sync.RWMutex
}

// NewHub 创建 Hub
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[uint64]*Room),
	}
}

func (rm *RoomManager) RemoverRoom(roomId uint64) {
	rm.mu.Lock()
	delete(rm.rooms, roomId)
	rm.mu.Unlock()
}

func (rm *RoomManager) GetRoom(roomId uint64) *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, ok := rm.rooms[roomId]
	if !ok {
		// 创建 room
		room = NewRoom(roomId)
		rm.rooms[roomId] = room
		go room.Run()
	}

	return room
}
