package redisService

import (
	"context"
	my_redis "ice-chat/pkg/redis"
	"time"
)

type roomRedisServiceImpl struct {
	redisOp my_redis.RedisOperator // 依赖接口，而非具体实现
}

type RoomRedisService interface {
	Allow(ctx context.Context, key string, ex time.Duration) bool
	GetRoomIDByJoinCode(ctx context.Context, key string) (uint64, error)
	CreateJoinCode(ctx context.Context, joinCode string, roomId uint64, ex time.Duration) error
}

func NewRoomRedisService(redisOp my_redis.RedisOperator) RoomRedisService {
	return &roomRedisServiceImpl{
		redisOp: redisOp,
	}
}

func (rs *roomRedisServiceImpl) Allow(ctx context.Context, key string, ex time.Duration) bool {
	ok, _ := rs.redisOp.SetNx(ctx, key, 1, ex)
	return ok
}

func (rs *roomRedisServiceImpl) GetRoomIDByJoinCode(ctx context.Context, key string) (uint64, error) {
	return rs.redisOp.GetUint64(ctx, key)
}

func (rs *roomRedisServiceImpl) CreateJoinCode(ctx context.Context, joinCode string, roomId uint64, ex time.Duration) error {
	return rs.redisOp.Set(ctx, joinCode, roomId, ex)
}
