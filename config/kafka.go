package config

type KafkaConfig struct {
	Brokers  []string          `mapstructure:"brokers"`
	Topic    string            `mapstructure:"topic"`
	GroupID  string            `mapstructure:"group_id"`
	Producer KafkaProducerConf `mapstructure:"producer"`
	Consumer KafkaConsumerConf `mapstructure:"consumer"`
}

type KafkaProducerConf struct {
	WriteTimeout int `mapstructure:"write_timeout"`
	ReadTimeout  int `mapstructure:"read_timeout"`
}

// KafkaConsumerConf Kafka消费者配置
type KafkaConsumerConf struct {
	MinBytes       int `mapstructure:"min_bytes"`
	MaxBytes       int `mapstructure:"max_bytes"`
	MaxWait        int `mapstructure:"max_wait"`
	ReadBackoffMin int `mapstructure:"read_backoff_min"`
}
