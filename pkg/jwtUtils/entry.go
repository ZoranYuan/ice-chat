package jwtUtils

// Init 初始化 jwt utils（必须调用）
func CreateJwtUtils(cfg Config) *jwtUtils {
	return newJWTUtils(cfg)
}
