package my_mysql

import (
	"ice-chat/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	client  *gorm.DB
	dbUtils *DBUtils // 包内初始化的工具实例
)

func Init() {
	// 配置 mysql 的 sql 输出
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 日志级别：Silent, Error, Warn, Info
			Colorful:      true,        // 彩色输出
		},
	)

	dsn := config.Conf.DB.DSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("❌ MySQL 初始化失败: %v", err)
	}
	client = db
	sqlDB, _ := client.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	dbUtils = NewDBUtils(client)
}

// GetDBUtils 获取工具实例（供业务层注入）
func GetDBUtils() *DBUtils {
	if dbUtils == nil {
		log.Panic("❌ 请先调用mysql.Init()")
	}
	return dbUtils
}
