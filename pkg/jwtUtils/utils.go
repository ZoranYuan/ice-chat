package jwtUtils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type jwtUtils struct {
	secret []byte
	expire time.Duration
}

func newJWTUtils(cfg Config) *jwtUtils {
	return &jwtUtils{
		secret: cfg.Secret,
		expire: cfg.Expire,
	}
}

func (j *jwtUtils) Generate(userId uint64) (string, *Claims, error) {
	jti := uuid.NewString()
	now := time.Now()

	claims := &Claims{
		UserID: userId,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(j.secret)
	if err != nil {
		return "", nil, err
	}

	return s, claims, nil
}

func (j *jwtUtils) Parse(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}
