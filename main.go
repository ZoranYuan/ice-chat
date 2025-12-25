package main

import (
	"ice-chat/config"
	"ice-chat/internal/api"
	"ice-chat/internal/redisService"
	"ice-chat/internal/repository"
	"ice-chat/internal/router"
	"ice-chat/internal/service"
	"ice-chat/pkg/mysql"
	"ice-chat/pkg/redis"
	"ice-chat/pkg/snowflake"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()
	redis.Init()
	mysql.Init()
	snowflake.Init()
	redisOp := redis.GetRedisOp()
	dbUtils := mysql.GetDBUtils() // db 只注入 resp 业务层中

	if err := dbUtils.AutoMigrate(); err != nil {
		log.Fatal(err)
	}

	userSvc := service.NewUserService(redisService.NewUserRepository(redisOp), repository.NewUserRepository(dbUtils))
	userAPI := api.NewUserAPI(userSvc)

	r := gin.Default()
	router.RegisterHeartCheckRouter(r)
	router.RegisterUserRouter(r, userAPI)

	if err := r.Run(config.Conf.App.GetAddress()); err != nil {
		log.Fatal(err)
	}
}
