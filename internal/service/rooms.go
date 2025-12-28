package service

import (
	"context"
	"errors"
	"ice-chat/internal/constants"
	model "ice-chat/internal/model/eneity"
	"ice-chat/internal/model/request"
	"ice-chat/internal/repository"
	"ice-chat/pkg/redis"
	"time"
)

type RoomsService interface {
	Create(group request.Room, uid uint64) (error, uint64)
}

type roomsServImpl struct {
	roomsRepo repository.RoomsRepository
	redisOp   redis.RedisOperator
}

func NewGroupsService(roomsRepo repository.RoomsRepository, redisOp redis.RedisOperator) RoomsService {
	return &roomsServImpl{
		roomsRepo: roomsRepo,
		redisOp:   redisOp,
	}
}

func (r *roomsServImpl) Create(room request.Room, uid uint64) (error, uint64) {
	ok, err := r.redisOp.SetNx(context.Background(), constants.USER_CREATE_GROUP, 60*time.Second)
	if err != nil {
		return err, 0
	}

	if !ok {
		return errors.New("操作频繁，请 60s 后再试"), 0
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
		return err, 0
	}

	return nil, roomsRepo.RoomID
}
