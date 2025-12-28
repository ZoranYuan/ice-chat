package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisOperator Redis操作接口（解耦核心）
type RedisOperator interface {
	Set(ctx context.Context, key string, value any, expire time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Del(ctx context.Context, key string) error
	SetNx(ctx context.Context, key string, expire time.Duration) (bool, error)
}

// RedisUtil 实现RedisOperator接口
type redisUtil struct {
	client *redis.Client
}

func NewRedisUtil(client *redis.Client) RedisOperator {
	return &redisUtil{client: client}
}

// Set 设置缓存（过期时间单位：秒）
func (ru *redisUtil) Set(ctx context.Context, key string, value any, expire time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis set marshal fail: %w", err)
	}
	return ru.client.Set(ctx, key, data, expire).Err()
}

// Get 获取缓存
func (ru *redisUtil) Get(ctx context.Context, key string, dest any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	data, err := ru.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("redis key [%s] not found", key)
		}
		return fmt.Errorf("redis get fail: %w", err)
	}
	return json.Unmarshal(data, dest)
}

// Del 删除缓存
func (ru *redisUtil) Del(ctx context.Context, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return ru.client.Del(ctx, key).Err()
}

// 唯一键设置
func (ru *redisUtil) SetNx(ctx context.Context, key string, expire time.Duration) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return ru.client.SetNX(ctx, key, "", expire).Result()
}
