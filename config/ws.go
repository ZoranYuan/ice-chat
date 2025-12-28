package config

import "time"

type WsConfig struct {
	ReadBufferSize  int `mapstructure:"read_buffer_size"`  // 读缓冲区大小（字节）
	WriteBufferSize int `mapstructure:"write_buffer_size"` // 写缓冲区大小（字节）

	PingIntervalSec  int `mapstructure:"ping_interval"`  // 心跳间隔（秒）- 原始配置值
	ReadDeadlineSec  int `mapstructure:"read_deadline"`  // 读超时（秒）- 原始配置值
	WriteDeadlineSec int `mapstructure:"write_deadline"` // 写超时（秒）- 原始配置值
	PingInterval     time.Duration
	ReadDeadline     time.Duration
	WriteDeadline    time.Duration
}

func (c *WsConfig) InitDuration() {
	c.PingInterval = time.Duration(c.PingIntervalSec) * time.Second
	c.ReadDeadline = time.Duration(c.ReadDeadlineSec) * time.Second
	c.WriteDeadline = time.Duration(c.WriteDeadlineSec) * time.Second
}
