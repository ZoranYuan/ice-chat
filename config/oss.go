package config

type OssConfig struct {
	EndPoint  string `mapstructure:"end_point"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Secure    bool   `mapstructure:"secure"`
}
