package snowflake

import (
	"time"

	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

// Init 初始化雪花ID生成器（程序启动时调用）
func Init() {
	// 设置起始时间（可选，建议设为项目启动时间）
	startTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	settings := sonyflake.Settings{
		StartTime: startTime,
	}
	sf = sonyflake.NewSonyflake(settings)
	if sf == nil {
		panic("初始化雪花ID生成器失败")
	}
}

// NewID 生成雪花ID（返回uint64类型，适配MySQL的bigint unsigned）
func NewID() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic("生成雪花ID失败: " + err.Error())
	}
	return id
}
