package service

import (
	"context"
	"encoding/json"
	"ice-chat/internal/model/request"
	"ice-chat/internal/mq/kafka"
	"ice-chat/internal/repository"
	"ice-chat/internal/ws"
	"log"
	"time"
)

type WsService struct {
	wsRepo    repository.MessageRepository
	kafka     *kafka.KafkaClient
	userRepo  repository.UserRepository
	wsManager *ws.ClientManager
}

const (
	// DB/Redis/Kafka操作超时
	opTimeout     = 5 * time.Second
	mainOpTimeout = 10 * time.Second
	// 消息重试间隔
	retryInterval = 10 * time.Second
	// 最大重试次数
	maxRetryCount = 3
	// Redis缓存失败消息的key前缀
	redisFailedKey = "ws:failed_msg:"
)

func NewWsService(wsRepo repository.MessageRepository, userRepo repository.UserRepository, kafka *kafka.KafkaClient, wsManager *ws.ClientManager) *WsService {
	return &WsService{
		wsRepo:    wsRepo,
		kafka:     kafka,
		userRepo:  userRepo,
		wsManager: wsManager,
	}
}

func (wss *WsService) MessageHandler(msgByte []byte, client *ws.Client) {
	_, mainCancel := context.WithTimeout(context.Background(), mainOpTimeout)
	defer mainCancel() // 主函数退出必执行，兜底防止ctx泄漏
	var message request.Message

	if err := json.Unmarshal(msgByte, &message); err != nil {
		log.Printf("fail to resolve ws message")

		// TODO 向客户端发送消息
		client.SendErrMessageToClient(err.Error())
		return
	}

	// 通过 kafka 来异步的完成 DB 操作
	// wss.kafka.Produce(context.Background(), msgByte)
	wss.wsManager.Broadcast(msgByte)
}

func (wss *WsService) handleKafkaProduce(ctx context.Context, msg []byte, client *ws.Client) {
	defer func() {
		if r := recover(); r != nil {
			// TODO 进一步处理
			log.Println(r)
			client.SendErrMessageToClient("发送失败")
		}
	}()

	if err := wss.kafka.Produce(ctx, msg); err != nil {
		panic(err)
	}
}
