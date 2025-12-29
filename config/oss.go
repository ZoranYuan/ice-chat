package config

type OssConfig struct {
	EndPoint        string `mapstructure:"end_point"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Secure          bool   `mapstructure:"secure"`
}
