package jwtUtils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	Secret []byte
	Expire time.Duration
}

type Claims struct {
	UserID uint64 `json:"userId"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}
