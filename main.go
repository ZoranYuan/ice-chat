package main

import (
	"ice-chat/config"
	"ice-chat/internal/api"
	"ice-chat/internal/mq/kafka"
	"ice-chat/internal/redisService"
	"ice-chat/internal/repository"
	"ice-chat/internal/router"
	"ice-chat/internal/service"
	"ice-chat/internal/ws"
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

	// 创建 ws 服务
	wsUtils := ws.GetUpgrader()
	wsChatManager := ws.NewManager()

	// 创建 mq
	kafkaClient := kafka.NewKafkaClient(wsChatManager)

	if err := dbUtils.AutoMigrate(); err != nil {
		log.Fatal(err)
	}

	// ws DI
	wsSvc := service.NewWsService(repository.NewUmsgRepository(dbUtils), repository.NewUserRepository(dbUtils), kafkaClient, wsChatManager)
	wsApi := api.NewWsAPI(wsSvc, wsChatManager, wsUtils)

	// user DI
	userSvc := service.NewUserService(redisService.NewUserRepository(redisOp), repository.NewUserRepository(dbUtils))
	userApi := api.NewUserAPI(userSvc)

	// group DI
	groupSvc := service.NewGroupsService(repository.NewGroupsRepo(dbUtils), redisOp)
	groupApi := api.NewGroupsApi(groupSvc)

	r := gin.Default()

	router.RegisterHeartCheckRouter(r)
	router.RegisterUserRouter(r, userApi)
	router.RegisterWsRouter(r, wsApi)
	router.RegisterGroupsRouter(r, groupApi)

	// 启动异步任务
	go wsChatManager.Run()
	go kafkaClient.Consume()

	if err := r.Run(config.Conf.App.GetAddress()); err != nil {
		log.Fatal(err)
	}
}
