package utils

import (
	"path/filepath"
	"strings"
)

var allowExts = map[string]bool{
	"mp4": true,
	"avi": true,
	"mov": true,
	"mkv": true,
	"flv": true,
	"wmv": true,
}

func getFileExt(fileName string) string {
	fileName = filepath.Base(fileName)
	dotIdx := strings.Index(fileName, ".")

	return fileName[dotIdx+1:]
}

func IsValidFileExt(fileName string) (string, bool) {
	ext := strings.ToLower(getFileExt(fileName)) // 统一转小写，兼容Video.MP4场景
	return ext, allowExts[ext]
}
