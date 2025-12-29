package service

import (
	"context"
	"fmt"
	"ice-chat/internal/constants"
	cerror "ice-chat/internal/model/error"
	"ice-chat/internal/model/request"
	"ice-chat/internal/repository"
	my_redis "ice-chat/pkg/redis"
	"ice-chat/scripts"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type UploadService interface {
	Merge(mergeInfo request.UploadMerge) error
	Upload(uploadId string, uploadChunkIdx int) error
	BreakPoint(uploadId string) (int64, error)
}
type uploadServiceImpl struct {
	userRepo repository.UserRepository
	minio    *minio.Client
	redisOp  my_redis.RedisOperator
}

func NewUploadService(userRepo repository.UserRepository, minio *minio.Client, redisOp my_redis.RedisOperator) UploadService {
	return &uploadServiceImpl{
		userRepo: userRepo,
		minio:    minio,
		redisOp:  redisOp,
	}
}

func (us *uploadServiceImpl) Merge(mergeInfo request.UploadMerge) error {
	// TODO 检索出路径下的所有文件，收集，写入到一个文件中，并且上传到 Minio 中
	tmpDir := fmt.Sprintf("%s/%s", constants.UPLOAD_TEMP_DIR, mergeInfo.UploadId)
	tmpFile := fmt.Sprintf("%s/%s.merge", constants.MERGE_TEMP_DIR, mergeInfo.UploadId)

	out, _ := os.Create(tmpFile)
	defer out.Close()

	for i := 0; i < mergeInfo.TotalChunks; i++ {
		part := fmt.Sprintf("%s/%d", tmpDir, i)
		in, err := os.Open(part)

		if err != nil {
			// TODO 让前端重新传递这个切片
			return &cerror.ChunkMissingError{MissingIndex: i}
		}

		if _, err := io.Copy(out, in); err != nil {
			in.Close()
			return err
		}
		in.Close()
	}

	// TODO 上传到 Minio
	f, _ := os.Open(tmpFile)
	stat, _ := f.Stat()

	objectKey := fmt.Sprintf("videos/%s.mp4", uuid.New().String())

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	if _, err := us.minio.PutObject(
		ctx,
		"videos",
		objectKey,
		f,
		stat.Size(),
		minio.PutObjectOptions{ContentType: "video/mp4"},
	); err != nil {
		return err
	}

	// TODO 将切片缓存从 redis 中删除
	defer cancel()
	idxKey := fmt.Sprintf("%s%s", constants.UPLOAD_CHUNK_IDX, mergeInfo.UploadId)
	err := us.redisOp.Del(ctx, idxKey)
	return err
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
	return us.redisOp.Set(ctx, idxKey, uploadChunkIdx, time.Duration(constants.UPLOAD_CHUNK_IDX_TIMEOUT))
}
