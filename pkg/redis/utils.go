package my_redis

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
	GetBytes(ctx context.Context, key string, dest any) error
	Del(ctx context.Context, key string) error
	GetUint64(ctx context.Context, key string) (uint64, error)
	SetNx(ctx context.Context, key string, value any, expire time.Duration) (bool, error)
	RunScript(ctx context.Context, redisScript *redis.Script, keys []string, args ...any) (int64, error)
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
func (ru *redisUtil) GetBytes(ctx context.Context, key string, dest any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	data, err := ru.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (ru *redisUtil) GetUint64(ctx context.Context, key string) (uint64, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	data, err := ru.client.Get(ctx, key).Uint64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, err
		}
		return 0, err
	}

	return data, nil
}

// Del 删除缓存
func (ru *redisUtil) Del(ctx context.Context, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return ru.client.Del(ctx, key).Err()
}

// 唯一键设置
func (ru *redisUtil) SetNx(ctx context.Context, key string, value any, expire time.Duration) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return ru.client.SetNX(ctx, key, value, expire).Result()
}

func (ru *redisUtil) RunScript(ctx context.Context, redisScript *redis.Script, keys []string, args ...any) (int64, error) {
	res, err := redisScript.Run(ctx, ru.client, keys, args...).Int64()

	return res, err
}
