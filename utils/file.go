package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
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
	dotIdx := strings.LastIndex(fileName, ".")
	return fileName[dotIdx+1:]
}

func IsValidFileExt(fileName string) (string, bool) {
	ext := strings.ToLower(getFileExt(fileName)) // 统一转小写，兼容Video.MP4场景
	return ext, allowExts[ext]
}

func TranscodeToMP4(inputPath, outputPath string) error {
	cmd := exec.Command(
		"D:\\software\\ffmpeg-8.0.1-full_build\\bin\\ffmpeg.exe",
		"-y",
		"-i", inputPath,
		"-c", "copy",
		outputPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %v, %s", err, stderr.String())
	}

	return nil
}

func GetVideoDuration(filePath string) (float64, error) {
	cmd := exec.Command(
		"D:\\software\\ffmpeg-8.0.1-full_build\\bin\\ffprobe.exe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return 0, err
	}

	durationStr := strings.TrimSpace(out.String())
	return strconv.ParseFloat(durationStr, 64)
}
