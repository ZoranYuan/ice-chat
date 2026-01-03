package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ice-chat/internal/constants"
	"ice-chat/internal/limter"
	"ice-chat/internal/model/request"
	res "ice-chat/internal/model/response"
	"ice-chat/internal/mq/kafka"
	"ice-chat/internal/repository"
	my_redis "ice-chat/pkg/redis"
	"ice-chat/pkg/ws"
	"ice-chat/scripts"
	"log"
	"strings"
	"time"
)

type wsServiceImpl struct {
	msgRepo   repository.MessageRepository
	kafka     *kafka.KafkaClient
	userRepo  repository.UserRepository
	wsManager *ws.RoomManager
	roomsRepo repository.RoomsRepository
	redisOp   my_redis.RedisOperator
	limter    *limter.Limter
}

type WsService interface {
	GroupIsExists(groupId uint64) (bool, error)
	ChatHandler(c *ws.Client, room *ws.Room, msg []byte)
	WatchHandler(c *ws.Client, room *ws.Room, msg []byte)
	SynchronizeVideoState(c *ws.Client, room *ws.Room)
}

func NewWsService(msgRepo repository.MessageRepository,
	userRepo repository.UserRepository,
	kafka *kafka.KafkaClient,
	wsManager *ws.RoomManager,
	roomsRepo repository.RoomsRepository,
	redisOp my_redis.RedisOperator,
) WsService {
	return &wsServiceImpl{
		msgRepo:   msgRepo,
		kafka:     kafka,
		userRepo:  userRepo,
		wsManager: wsManager,
		roomsRepo: roomsRepo,
		redisOp:   redisOp,
		limter:    limter.NewLimter(5, 10),
	}
}

func (wss *wsServiceImpl) GroupIsExists(groupId uint64) (bool, error) {
	return wss.roomsRepo.RoomIsExists(groupId)
}

func (wss *wsServiceImpl) handle(c *ws.Client, room *ws.Room, msg []byte) {
	_, mainCancel := context.WithTimeout(context.Background(), time.Duration(constants.WS_TIMEOUT_EXPIRETIME)*time.Second)
	defer mainCancel()
	// TODO 这里暂时只做一个广播的作用
	room.Broadcast(msg)
}

func (wss *wsServiceImpl) ChatHandler(c *ws.Client, room *ws.Room, msg []byte) {
	var message request.Message

	// 单机限流
	if allow := wss.limter.Allow(c.GetUid()); !allow {
		return
	}

	if err := json.Unmarshal(msg, &message); err != nil {
		log.Printf("非法消息: %v", err)
		ws.Fail(c, "消息格式错误")
		return
	}

	// 可选：进一步校验字段
	content := strings.TrimSpace(message.Content)
	if content == "" {
		return
	}

	if len(content) > 500 {
		ws.Fail(c, "消息过长")
		return
	}

	wss.handle(c, room, msg)
}

func (wss *wsServiceImpl) WatchHandler(c *ws.Client, room *ws.Room, msg []byte) {
	// TODO 将状态同步到 redis 中
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(constants.REDIS_TIMEOUT)*time.Millisecond)
	defer cancel()

	key := fmt.Sprintf("%s%d:%d", constants.VIDEO_CONTROL_LOCK, c.GetUid(), room.GetRoomId())
	log.Println(key)
	// 限流锁, 这里前端会做防抖处理
	if ok, _ := wss.redisOp.SetNx(ctx, key, "", time.Duration(constants.VIDEO_CONTROL_INTERVAL)*time.Second); !ok {
		c.SendMessageToClient([]byte("操作太频繁啦"))
		return
	}
	var message request.VideoState
	decoder := json.NewDecoder(bytes.NewReader(msg))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&message); err != nil {
		log.Printf("failed to marsha1: %v", err)
		ws.Fail(c, "非法请求")
		return
	}

	if message.TimeStamp <= 0 {
		// 非法参数
		ws.Fail(c, "非法请求")
		return
	}

	// TODO 原子更新状态
	tsKey := fmt.Sprintf("%s%d", constants.VIDEO_LATEST_TIMESTAMP, room.GetRoomId())
	stateKey := fmt.Sprintf("%s%d", constants.VIDEO_ROOM_STATE, room.GetRoomId())
	res, err := wss.redisOp.RunScript(ctx, scripts.UpdateLatestTimestamp, []string{tsKey, stateKey}, constants.VIDEO_STATE_TIME, message.TimeStamp, string(msg))
	if err != nil {
		log.Printf("lua script failed: %v", err)
		ws.Fail(c, "服务异常")
		return
	}

	if res == 0 {
		// 旧指令，正常丢弃
		return
	}

	// 广播
	wss.handle(c, room, msg)
}

func (wss *wsServiceImpl) SynchronizeVideoState(c *ws.Client, room *ws.Room) {
	var videoState res.VideoStateInit
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	stateKey := fmt.Sprintf("%s%d", constants.VIDEO_STATE_ROOM_INIT, room.GetRoomId())
	for {
		if err := wss.redisOp.GetBytes(ctx, stateKey, &videoState); err != nil {
			log.Println("debug")
			time.Sleep(1000 * time.Millisecond)
			continue
		} else {
			break
		}
	}

	ws.Ok(c, videoState)
}
