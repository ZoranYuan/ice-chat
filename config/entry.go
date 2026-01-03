package config

import (
	"flag"
	"log"

	"github.com/spf13/viper"
)

var Conf *config

func Init() {
	// TODO 对 config 结构体进行初始化
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	viper.SetConfigFile(".env")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("❌ 读取配置文件失败: %v", err)
	}
	if err := v.Unmarshal(&Conf); err != nil {
		log.Fatalf("❌ 解析配置文件失败: %v", err)
	}

	// 初始化时间
	Conf.Ws.InitDuration()

	// 获取命令行的 port 参数
	var port = flag.Int("port", Conf.App.Port, "server port")
	flag.Parse()
	if port != nil {
		Conf.App.Port = *port
	} else {
		Conf.App.Port = 8081
	}
}
