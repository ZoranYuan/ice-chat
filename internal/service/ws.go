package service

import (
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

	// 限流锁, 这里前端会做防抖处理
	key := fmt.Sprintf("%s%d:%d", constants.VIDEO_CONTROL_LOCK, c.GetUid(), room.GetRoomId())
	if ok, _ := wss.redisOp.SetNx(ctx, key, "", time.Duration(constants.VIDEO_CONTROL_INTERVAL)*time.Millisecond); !ok {
		c.SendMessageToClient([]byte("操作太频繁啦"))
		return
	}

	// TODO 引入并发锁， 对于 A 和 B 两个请求来说，如果是同时发生的，可能出现覆盖的问题，比如 A 来更新（读取到旧状态），此时旧状态未能及时更新，B 又来（读取到旧状态）
	/*
		前端传递的 msg 格式

		type VideoCommand struct {
			RoomId    uint64
			Action    VideoAction
			Time      float64      // seek 时使用，前端拖拽进度条的时候使用
			Speed     float64
			TimeStamp int64      // 防止旧消息
		}
	*/
	now := time.Now().Unix()

	var videoCommand request.VideoCommand
	if err := json.Unmarshal(msg, &videoCommand); err != nil {
		log.Println("Failed to get video state", err)
		return
	}

	videoStateKey := fmt.Sprintf("%s%d", constants.VIDEO_ROOM_STATE, videoCommand.RoomId)

	// TODO 直接使用 Lua 脚本来避免并发带来的覆盖问题
	result, err := wss.redisOp.RunScriptWithData(ctx, scripts.UpdateVideoState, []string{videoStateKey}, now, videoCommand.Action,
		videoCommand.Action,
		videoCommand.Time,
		videoCommand.Speed,
		videoCommand.TimeStamp)
	if err != nil {
		log.Println("Failed to resolve, ", err)
		return
	}

	if result == nil {
		log.Println("Failed to update video state")
		return
	}

	newStateJson := result.(string)
	// TODO 广播
	wss.handle(c, room, []byte(newStateJson))
}

func (wss *wsServiceImpl) SynchronizeVideoState(c *ws.Client, room *ws.Room) {
	var videoState res.VideoState
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	videoStateKey := fmt.Sprintf("%s%d", constants.VIDEO_ROOM_STATE, room.GetRoomId())
	if err := wss.redisOp.GetBytes(ctx, videoStateKey, &videoState); err != nil {
		log.Println("debug")
		ws.Fail(c, "视频同步失败，请稍后再试")
	} else {
		ws.Ok(c, videoState)
	}
}
