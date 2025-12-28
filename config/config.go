package config

type config struct {
	App   AppConfig   `mapstructure:"app"`
	DB    DBConfig    `mapstructure:"db"`
	Redis RedisConfig `mapstructure:"redis"`
	JWT   JWTConfig   `mapstructure:"jwt"`
	Ws    WsConfig    `mapstructure:"ws"`
	Kafka KafkaConfig `mapstructure:"kafka"`
}
