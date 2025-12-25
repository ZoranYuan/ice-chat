package redis

import (
	"context"
	"ice-chat/config"
	"log"

	"github.com/redis/go-redis/v9"
)

// 包内私有变量：客户端+工具实例（对外不暴露）
var (
	client  *redis.Client
	redisOp RedisOperator // 包内初始化的工具实例
)

// Init 初始化Redis（对外唯一的初始化入口）
func Init() {
	redisConf := config.Conf.Redis

	// 1. 创建Redis客户端
	client = redis.NewClient(&redis.Options{
		Addr:     redisConf.GetAddress(),
		Password: redisConf.Password,
		DB:       redisConf.DB,
	})

	// 2. 测试连接
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("❌ Redis初始化失败: %v", err)
	}

	// 3. 初始化Redis工具实例（包内私有）
	redisOp = NewRedisUtil(client)

	log.Println("✅ Redis初始化成功")
}

// GetUtil 获取Redis工具实例（对外暴露，供业务层注入）
func GetRedisOp() RedisOperator {
	if redisOp == nil {
		log.Panic("❌ 请先调用redis.Init()初始化Redis")
	}
	return redisOp
}
