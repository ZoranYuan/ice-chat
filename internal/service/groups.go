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

type GroupsService interface {
	Create(group request.Group, uid uint64) (error, uint64)
}

type groupServImpl struct {
	groupRepo repository.GroupsRepository
	redisOp   redis.RedisOperator
}

func NewGroupsService(groupRepo repository.GroupsRepository, redisOp redis.RedisOperator) GroupsService {
	return &groupServImpl{
		groupRepo: groupRepo,
		redisOp:   redisOp,
	}
}

func (g *groupServImpl) Create(group request.Group, uid uint64) (error, uint64) {
	ok, err := g.redisOp.SetNx(context.Background(), constants.USER_CREATE_GROUP, 60*time.Second)
	if err != nil {
		return err, 0
	}

	if !ok {
		return errors.New("操作频繁，请 60s 后再试"), 0
	}

	groupRepo := model.Groups{
		GroupName:  group.GroupName,
		Avatar:     group.Avatar,
		Desc:       group.Desc,
		CreateUser: uid,
	}

	groupMemberRepo := model.GroupMember{
		UserID: uid,
	}

	if err := g.groupRepo.Create(&groupRepo, &groupMemberRepo); err != nil {
		return err, 0
	}

	return nil, groupRepo.GroupId
}
