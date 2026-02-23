package redisService

import (
	"context"
	"ice-chat/config"
	my_redis "ice-chat/pkg/redis"
)

type userRedisService struct {
	redisOp my_redis.RedisOperator // 依赖接口，而非具体实现
}

type UserReids interface {
	StoreAccessKey(key, value string) error
	DeleteAccessKey(key string) error
}

func NewUserRepository(redisOp my_redis.RedisOperator) UserReids {
	return &userRedisService{redisOp: redisOp}
}

func (u *userRedisService) StoreAccessKey(key, value string) error {
	err := u.redisOp.Set(context.TODO(), key, value, config.Conf.JWT.GetAccessTokenExpireDuration())
	return err
}

func (u *userRedisService) DeleteAccessKey(key string) error {
	err := u.redisOp.Del(context.TODO(), key)
	return err
}
