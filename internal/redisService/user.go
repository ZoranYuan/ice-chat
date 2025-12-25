package redisService

import (
	"context"
	"ice-chat/config"
	"ice-chat/pkg/redis"
)

type userRedisService struct {
	redisOp redis.RedisOperator // 依赖接口，而非具体实现
}

type UserReids interface {
	StoreAccessKey(key string) error
}

func NewUserRepository(redisOp redis.RedisOperator) UserReids {
	return &userRedisService{redisOp: redisOp}
}

func (u *userRedisService) StoreAccessKey(key string) error {
	err := u.redisOp.Set(context.TODO(), key, struct{}{}, config.Conf.JWT.GetAccessTokenExpireDuration())
	return err
}
