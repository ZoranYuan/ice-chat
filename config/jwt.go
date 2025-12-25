package config

import (
	"log"
	"time"
)

type JWTConfig struct {
	Secret                     string `mapstructure:"secret"`
	AccessTokenExpireDuration  string `mapstructure:"access_token_expire_duration"`
	RefreshTokenExpireDuration string `mapstructure:"refresh_token_expire_duration"`
}

func (j *JWTConfig) GetAccessTokenExpireDuration() time.Duration {
	dur, err := time.ParseDuration(j.AccessTokenExpireDuration)
	if err != nil {
		log.Fatal(err)
	}
	return dur
}

func (j *JWTConfig) GetRefreshTokenExpireDuration() time.Duration {
	dur, err := time.ParseDuration(j.RefreshTokenExpireDuration)
	if err != nil {
		log.Fatal(err)
	}
	return dur
}
