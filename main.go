package main

import (
	"ice-chat/config"
	"ice-chat/internal/api"
	"ice-chat/internal/mq/kafka"
	"ice-chat/internal/oss"
	"ice-chat/internal/redisService"
	"ice-chat/internal/repository"
	"ice-chat/internal/router"
	"ice-chat/internal/service"
	my_mysql "ice-chat/pkg/mysql"
	my_redis "ice-chat/pkg/redis"
	"ice-chat/pkg/snowflake"
	"ice-chat/pkg/ws"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()
	my_redis.Init()
	my_mysql.Init()
	snowflake.Init()
	redisOp := my_redis.GetRedisOp()
	dbUtils := my_mysql.GetDBUtils() // db 只注入 resp 业务层中
	minioClient := oss.NewMinioClient(config.Conf.Oss)
	// 创建 ws 服务
	wsUtils := ws.NewWsUtils()
	roomManager := ws.NewRoomManager()
	// TODO 版本更新后 ，kafka 功能待完善
	kafkaClient := kafka.NewKafkaClient(roomManager)

	if err := dbUtils.AutoMigrate(); err != nil {
		log.Fatal(err)
	}

	// chat DI
	wsSvc := service.NewWsService(repository.NewUmsgRepository(dbUtils), repository.NewUserRepository(dbUtils), kafkaClient, roomManager, repository.NewRoomsRepo(dbUtils), redisOp)
	wsApi := api.NewWsAPI(wsSvc, roomManager, wsUtils)

	// user DI
	userSvc := service.NewUserService(redisService.NewUserRepository(redisOp), repository.NewUserRepository(dbUtils), minioClient)
	userApi := api.NewUserAPI(userSvc)

	// room DI
	roomRedisService := redisService.NewRoomRedisService(redisOp)
	roomSev := service.NewRoomService(repository.NewRoomsRepo(dbUtils), roomRedisService)
	roomApi := api.NewRoomsApi(roomSev)

	// upload
	uploadSev := service.NewUploadService(repository.NewUserRepository(dbUtils), minioClient, redisOp)
	uploadApi := api.NewUploadApi(uploadSev)

	r := gin.Default()

	router.RegisterHeartCheckRouter(r)
	router.RegisterUserRouter(r, userApi)
	router.RegisterWsRouter(r, wsApi)
	router.RegisterRoomsRouter(r, roomApi)
	router.RegisterUploadRouter(r, uploadApi)
	// 启动异步任务
	// go kafkaClient.Consume()

	if err := r.Run(config.Conf.App.GetAddress()); err != nil {
		log.Fatal(err)
	}
}
