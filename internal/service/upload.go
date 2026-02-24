package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ice-chat/internal/constants"
	"ice-chat/internal/model/request"
	res "ice-chat/internal/model/response"
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
	GetMergeState(uploadId string) (uint64, error)
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

// TODO 1. 支持断点续传 （通过 redis 存储已经上传的切片的索引 -- 注意：上传失败的文件，前端已经存储，前端会对所有未成功上传的切片进行重新传递）
// TODO 2. 合并切片的时候，实时更新 redis 状态 (交由前端获取最新的状态)，如果获取到的状态值为 -1，那么表示 merge 失败了，此时优化点就是：后端会在数据库中对已经上传的文件进行落盘（存储已经合并了哪些文件，最后由前端请求，再重新在断点进行合并）

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
	go func() {
		mergeStateKey := fmt.Sprintf("%s%s", constants.VIDEO_MERGE_KEY, mergeInfo.UploadId)

		mergeI := 1

		defer func() {
			// TODO 切片合并失败
			mergeState := res.MergeState{
				Merged: mergeI,
				Status: res.MergeFailed,
				Total:  mergeInfo.TotalChunks,
				Error:  "",
			}

			data, err := json.Marshal(mergeState)
			if err != nil {
				log.Println("Failed to marshal ", err)
			}

			if err := us.redisOp.Set(ctx, mergeStateKey, data, 10*time.Minute); err != nil {
				return
			}
		}()

		for ; mergeI <= mergeInfo.TotalChunks; mergeI++ {
			part := fmt.Sprintf("%s/%d", tempDir, mergeI)
			in, err := os.Open(part)
			if err != nil {
				return
			}

			_, err = io.Copy(out, in)
			in.Close()
			if err != nil {
				return
			}

			// TODO 记录合并的状态， 可以将状态值落盘在数据库中
			mergeState := res.MergeState{
				Merged: mergeI,
				Status: res.MergeMerging,
				Total:  mergeInfo.TotalChunks,
				Error:  "",
			}

			data, err := json.Marshal(mergeState)
			if err != nil {
				log.Println("Failed to marshal ", err)
			}

			if err := us.redisOp.Set(ctx, mergeStateKey, data, 10*time.Minute); err != nil {
				return
			}
		}

		mergeState := res.MergeState{
			Merged: mergeI,
			Status: res.MergeSuccess,
			Total:  mergeInfo.TotalChunks,
			Error:  "",
		}

		data, err := json.Marshal(mergeState)
		if err != nil {
			log.Println("Failed to marshal ", err)
		}

		if err := us.redisOp.Set(ctx, mergeStateKey, data, 10*time.Minute); err != nil {
			return
		}
	}()

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
	result, err := us.redisOp.RunScript(ctx, scripts.UpdateUploadChunkIdx, []string{idxKey}, time.Duration(constants.UPLOAD_CHUNK_IDX_TIMEOUT)*time.Hour)
	return result, err
}

func (us *uploadServiceImpl) Upload(uploadId string, uploadChunkIdx int) error {
	// TODO 将当前发送的 chunk 更新到 Redis 中
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(constants.REDIS_TIMEOUT)*time.Millisecond)
	defer cancel()
	idxKey := fmt.Sprintf("%s%s", constants.UPLOAD_CHUNK_IDX, uploadId)
	return us.redisOp.Set(ctx, idxKey, uploadChunkIdx, time.Duration(constants.UPLOAD_CHUNK_IDX_TIMEOUT)*time.Hour)
}

func (us *uploadServiceImpl) GetMergeState(uploadId string) (uint64, error) {
	mergeStateKey := fmt.Sprintf("%s%s", constants.VIDEO_MERGE_KEY, uploadId)
	return us.redisOp.GetUint64(context.Background(), mergeStateKey)
}
