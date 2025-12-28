package service

import (
	"context"
	"encoding/json"
	"ice-chat/internal/constants"
	"ice-chat/internal/model/request"
	"ice-chat/internal/mq/kafka"
	"ice-chat/internal/repository"
	"ice-chat/pkg/ws"
	"log"
	"time"
)

type wsServiceImpl struct {
	msgRepo   repository.MessageRepository
	kafka     *kafka.KafkaClient
	userRepo  repository.UserRepository
	wsManager *ws.RoomManager
	roomsRepo repository.RoomsRepository
}

type WsService interface {
	GroupIsExists(groupId uint64) (bool, error)
	ChatHandler(c *ws.Client, room *ws.Room, msg []byte)
	WatchHandler(c *ws.Client, room *ws.Room, msg []byte)
}

func NewWsService(msgRepo repository.MessageRepository, userRepo repository.UserRepository, kafka *kafka.KafkaClient, wsManager *ws.RoomManager, roomsRepo repository.RoomsRepository) WsService {
	return &wsServiceImpl{
		msgRepo:   msgRepo,
		kafka:     kafka,
		userRepo:  userRepo,
		wsManager: wsManager,
		roomsRepo: roomsRepo,
	}
}

func (wss *wsServiceImpl) GroupIsExists(groupId uint64) (bool, error) {
	return wss.roomsRepo.GroupIsExists(groupId)
}

func (wss *wsServiceImpl) handle(c *ws.Client, room *ws.Room, msg []byte) {
	_, mainCancel := context.WithTimeout(context.Background(), time.Duration(constants.WS_EXPIRETIME)*time.Second)
	defer mainCancel()

	// TODO 这里暂时只做一个广播的作用,需要将消息存储到 db 中
	var message request.Message

	if err := json.Unmarshal(msg, &message); err != nil {
		// 消息不符合 Message 结构
		log.Printf("非法消息: %v", err)
		room.RemoveClient(c)
		return
	}

	// 可选：进一步校验字段
	if message.Content == "" {
		log.Println("消息内容为空，忽略")
		return
	}

	room.Broadcast(msg)
}

func (wss *wsServiceImpl) ChatHandler(c *ws.Client, room *ws.Room, msg []byte) {
	wss.handle(c, room, msg)
}

func (wss *wsServiceImpl) WatchHandler(c *ws.Client, room *ws.Room, msg []byte) {
	wss.handle(c, room, msg)

}
