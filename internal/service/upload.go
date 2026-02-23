package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ice-chat/internal/constants"
	cerror "ice-chat/internal/model/error"
	"ice-chat/internal/model/request"
	mq_eneity "ice-chat/internal/mq/entity"
	"ice-chat/internal/mq/kafka"
	"ice-chat/internal/repository"
	my_redis "ice-chat/pkg/redis"
	"ice-chat/scripts"
	"io"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
)

type UploadService interface {
	Merge(mergeInfo request.UploadMerge) error
	Upload(uploadId string, uploadChunkIdx int) error
	BreakPoint(uploadId string) (int64, error)
}

type uploadServiceImpl struct {
	userRepo    repository.UserRepository
	minio       *minio.Client
	redisOp     my_redis.RedisOperator
	kafkaClient *kafka.KafkaClient
}

func NewUploadService(userRepo repository.UserRepository, minio *minio.Client, redisOp my_redis.RedisOperator, kafkaClient *kafka.KafkaClient) UploadService {
	return &uploadServiceImpl{
		userRepo:    userRepo,
		minio:       minio,
		redisOp:     redisOp,
		kafkaClient: kafkaClient,
	}
}

func (us *uploadServiceImpl) Merge(mergeInfo request.UploadMerge) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tempDir := fmt.Sprintf("%s/%s", constants.UPLOAD_TEMP_DIR, mergeInfo.UploadId)
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return fmt.Errorf("upload temp dir does not exist: %s", tempDir)
	}

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil && !os.IsNotExist(err) {
			log.Println("Failed to remove temp file:", tempDir, err)
		}
	}()

	// 确保合并目录存在
	if err := os.MkdirAll(constants.MERGE_TEMP_DIR, 0755); err != nil {
		return err
	}

	tmpFile := fmt.Sprintf("%s/%s.tmp", constants.MERGE_TEMP_DIR, mergeInfo.UploadId)
	outFile := fmt.Sprintf("%s/%s.mp4", constants.MERGE_TEMP_DIR, mergeInfo.UploadId)

	// 创建临时文件并检查错误
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create tmp file: %w", err)
	}
	defer out.Close()

	// 合并所有 chunk
	for i := 0; i < mergeInfo.TotalChunks; i++ {
		part := fmt.Sprintf("%s/%d", tempDir, i)
		in, err := os.Open(part)
		if err != nil {
			return &cerror.ChunkMissingError{MissingIndex: i}
		}

		n, err := io.Copy(out, in)
		in.Close()
		if err != nil {
			return fmt.Errorf("failed to copy chunk %d: %w", i, err)
		}
		if n == 0 {
			return fmt.Errorf("chunk %d is empty", i)
		}
	}

	// 检查合并后文件大小
	info, err := out.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat tmp file: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("merged file is empty")
	}

	// 发送转码任务
	transcodeTask := mq_eneity.TranscodeTask{
		TmepFile: tmpFile,
		OutFile:  outFile,
		RoomID:   mergeInfo.RoomId,
		UploadID: mergeInfo.UploadId,
	}
	msg, err := json.Marshal(transcodeTask)
	if err != nil {
		return err
	}
	if err := us.kafkaClient.Produce(ctx, msg, "video-transcode"); err != nil {
		return err
	}

	return nil
}

func (us *uploadServiceImpl) BreakPoint(uploadId string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(constants.REDIS_TIMEOUT)*time.Millisecond)
	defer cancel()
	idxKey := fmt.Sprintf("%s%s", constants.UPLOAD_CHUNK_IDX, uploadId)
	res, err := us.redisOp.RunScript(ctx, scripts.UpdateUploadChunkIdx, []string{idxKey}, time.Duration(constants.UPLOAD_CHUNK_IDX_TIMEOUT)*time.Hour)
	return res, err
}

func (us *uploadServiceImpl) Upload(uploadId string, uploadChunkIdx int) error {
	// TODO 将当前发送的 chunk 更新到 Redis 中
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(constants.REDIS_TIMEOUT)*time.Millisecond)
	defer cancel()
	idxKey := fmt.Sprintf("%s%s", constants.UPLOAD_CHUNK_IDX, uploadId)
	return us.redisOp.Set(ctx, idxKey, uploadChunkIdx, time.Duration(constants.UPLOAD_CHUNK_IDX_TIMEOUT)*time.Hour)
}
