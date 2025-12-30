package service

import (
	"context"
	"errors"
	"fmt"
	"ice-chat/internal/constants"
	model "ice-chat/internal/model/eneity"
	"ice-chat/internal/model/request"
	"ice-chat/internal/redisService"
	"ice-chat/internal/repository"
	"ice-chat/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type RoomsService interface {
	CreateRoom(group request.Room, uid uint64) (error, uint64, string)
	JoinRoom(uid uint64, joinCode string) (uint64, error)
}

type roomsServImpl struct {
	roomsRepo        repository.RoomsRepository
	roomRedisService redisService.RoomRedisService
}

func NewRoomService(roomsRepo repository.RoomsRepository, roomRedisService redisService.RoomRedisService) RoomsService {
	return &roomsServImpl{
		roomsRepo:        roomsRepo,
		roomRedisService: roomRedisService,
	}
}

/*
@desc 单用户限流
*/
func (r *roomsServImpl) Allow(limitKey string, ex time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(constants.REDIS_TIMEOUT)*time.Millisecond)
	defer cancel()

	return r.roomRedisService.Allow(ctx, limitKey, ex)
}

func (r *roomsServImpl) CreateRoom(room request.Room, uid uint64) (error, uint64, string) {
	limitKey := fmt.Sprintf("%s%d", constants.USER_CREATE_GROUP_LOCK, uid)
	if ok := r.Allow(limitKey, time.Duration(constants.USER_CREATE_GROUP_TIMEINTERVAL)*time.Millisecond); !ok {
		return errors.New("操作频繁"), 0, ""
	}

	roomsRepo := model.Rooms{
		RoomName:   room.RoomName,
		Avatar:     room.Avatar,
		Desc:       room.Desc,
		CreateUser: uid,
	}

	groupMemberRepo := model.RoomsMember{
		UserID: uid,
	}

	if err := r.roomsRepo.Create(&roomsRepo, &groupMemberRepo); err != nil {
		return err, 0, ""
	}

	// 产生邀请码
	joinCode := utils.GenJoinCode(6)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(constants.REDIS_TIMEOUT)*time.Millisecond)
	defer cancel()
	joinCodeKey := fmt.Sprintf("%s%s", constants.ROOM_JOINCODE, joinCode)
	if err := r.roomRedisService.CreateJoinCode(ctx, joinCodeKey, roomsRepo.RoomID, time.Duration(constants.ROOM_JOINCODE_EFFECTIVE_TIME)*time.Millisecond); err != nil {
		return err, 0, ""
	}

	return nil, roomsRepo.RoomID, joinCode
}

func (r *roomsServImpl) JoinRoom(uid uint64, joinCode string) (uint64, error) {
	limitKey := fmt.Sprintf("%s%d", constants.ROOM_JOIN_UID_LOCK, uid)

	if ok := r.Allow(limitKey, time.Duration(constants.ROOM_JOIN_USER_TIMEINTERVAL)*time.Millisecond); !ok {
		return 0, errors.New("操作频繁")
	}

	// TODO 查询对应的 room
	joinCodeKey := fmt.Sprintf("%s%s", constants.ROOM_JOINCODE, joinCode)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(constants.REDIS_TIMEOUT)*time.Millisecond)
	defer cancel()
	roomId, err := r.roomRedisService.GetRoomIDByJoinCode(ctx, joinCodeKey)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			// TODO 优化方向：使用 redis hash 来存储： key 就是 joinCode，value 为 roomId 和 过期时间戳
			return 0, errors.New("无效的验证码")
		}

		return 0, err
	}

	// TODO 查询 DB （后期可以在 redis 中优化）
	exists, err := r.roomsRepo.RoomIsExists(roomId)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, errors.New("房间不存在")
	}

	// 开启事务，将当前用户加入到对应的 roomId 中
	if err := r.roomsRepo.JoinRoom(uid, roomId); err != nil {
		return 0, err
	}

	return roomId, nil
}
