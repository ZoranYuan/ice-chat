package jwtUtils

func CreateJwtUtils(cfg Config) *jwtUtils {
	return newJWTUtils(cfg)
}
