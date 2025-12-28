package main

import (
	"ice-chat/config"
	"ice-chat/internal/api"
	"ice-chat/internal/mq/kafka"
	"ice-chat/internal/redisService"
	"ice-chat/internal/repository"
	"ice-chat/internal/router"
	"ice-chat/internal/service"
	"ice-chat/pkg/mysql"
	"ice-chat/pkg/redis"
	"ice-chat/pkg/snowflake"
	"ice-chat/pkg/ws"
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
	wsUtils := ws.NewWsUtils()
	roomManager := ws.NewRoomManager()
	// TODO 版本更新后 ，kafka 功能待完善
	kafkaClient := kafka.NewKafkaClient(roomManager)

	if err := dbUtils.AutoMigrate(); err != nil {
		log.Fatal(err)
	}

	// chat DI
	wsSvc := service.NewWsService(repository.NewUmsgRepository(dbUtils), repository.NewUserRepository(dbUtils), kafkaClient, roomManager, repository.NewGroupsRepo(dbUtils))
	wsApi := api.NewWsAPI(wsSvc, roomManager, wsUtils)

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
	// go kafkaClient.Consume()

	if err := r.Run(config.Conf.App.GetAddress()); err != nil {
		log.Fatal(err)
	}
}
