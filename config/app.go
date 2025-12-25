package config

import "fmt"

type AppConfig struct {
	Name string `mapstructure:"name"`
	Mode string `mapstructure:"mode"`
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

func (appConfig *AppConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port)
}
