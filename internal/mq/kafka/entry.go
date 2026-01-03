package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"ice-chat/config"
	"ice-chat/internal/constants"
	res "ice-chat/internal/model/response"
	mq_eneity "ice-chat/internal/mq/entity"
	my_redis "ice-chat/pkg/redis"
	"ice-chat/utils"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/segmentio/kafka-go"
)

// KafkaClient =
type KafkaClient struct {
	writer  *kafka.Writer
	readers map[string]*kafka.Reader
	minio   *minio.Client
	redisOp my_redis.RedisOperator
}

var topicHandlers map[string]func(kafka.Message) error

func NewKafkaClient(topics []string, minio *minio.Client, redisOp my_redis.RedisOperator) *KafkaClient {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Conf.Kafka.Brokers...),
		Balancer:               &kafka.LeastBytes{},
		WriteTimeout:           time.Duration(config.Conf.Kafka.Producer.WriteTimeout),
		BatchSize:              1,
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireOne,
		Async:                  true,
	}

	readers := make(map[string]*kafka.Reader)
	for _, topic := range topics {
		readers[topic] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:        config.Conf.Kafka.Brokers,
			Topic:          topic,
			GroupID:        config.Conf.Kafka.GroupID,
			MinBytes:       config.Conf.Kafka.Consumer.MinBytes,
			MaxBytes:       config.Conf.Kafka.Consumer.MaxBytes,
			ReadBackoffMin: time.Duration(config.Conf.Kafka.Consumer.ReadBackoffMin),
			Dialer: &kafka.Dialer{
				Timeout:   120 * time.Second,
				KeepAlive: 60 * time.Second,
			},
			StartOffset:    kafka.LastOffset,
			CommitInterval: 0,
		})
	}

	return &KafkaClient{
		writer:  writer,
		readers: readers,
		minio:   minio,
		redisOp: redisOp,
	}
}

func (k *KafkaClient) StartConsumers(topic []string) {
	topicHandlers = make(map[string]func(kafka.Message) error)
	topicHandlers["video-transcode"] = k.handleVideoTranscode

	for topic, reader := range k.readers {
		handler, ok := topicHandlers[topic]
		if !ok {
			log.Printf("No handler for topic: %s", topic)
			continue
		}
		go func(r *kafka.Reader, h func(kafka.Message) error, t string) {
			defer r.Close()
			for {
				msg, err := r.ReadMessage(context.Background())
				if err != nil {
					log.Printf("[%s] Kafka read error: %v", t, err)
					time.Sleep(time.Second)
					continue
				}
				if err := h(msg); err != nil {
					log.Printf("[%s] handle message failed: %v", t, err)
					// 不提交 offset，让 Kafka 重试
					continue
				}

				if err := r.CommitMessages(context.Background(), msg); err != nil {
					log.Println("Failed to commit offset:", err)
				}
			}
		}(reader, handler, topic)
	}
}

func (k *KafkaClient) Produce(ctx context.Context, msg []byte, topic string) error {
	return k.writer.WriteMessages(ctx,
		kafka.Message{
			Value: msg,
			Topic: topic,
		},
	)
}

func (k *KafkaClient) handleVideoTranscode(msg kafka.Message) error {
	var task mq_eneity.TranscodeTask

	if ok := utils.UnmarshalTools(msg.Value, &task); !ok {
		log.Println("Failed to unmarshal task")
		return fmt.Errorf("unmarshal failed")
	}

	if ok, err := k.redisOp.SetNx(context.TODO(), constants.VIDEO_TRANSCODE_KEY+task.UploadID, "", time.Duration(constants.VIDEOS_TRANSCODE_TIMEOUT)*time.Millisecond); !ok {
		return err
	}

	// 转码
	if err := utils.TranscodeToMP4(task.TmepFile, task.OutFile); err != nil {
		log.Println("Failed to transcode:", err)
		return err
	}

	// 打开文件
	f, err := os.Open(task.OutFile)
	if err != nil {
		log.Println("Failed to open outFile:", err)
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		log.Println("Failed to stat outFile:", err)
		return err
	}
	if stat.Size() == 0 {
		log.Println("OutFile is empty, skipping upload")
		return fmt.Errorf("empty file")
	}

	// 上传到 MinIO
	objectKey := task.UploadID + ".mp4"
	if _, err := k.minio.PutObject(context.Background(),
		"videos",
		objectKey,
		f,
		stat.Size(),
		minio.PutObjectOptions{ContentType: "video/mp4"},
	); err != nil {
		log.Println("Upload to Minio failed:", err)
		return err
	}

	videoDuration, err := utils.GetVideoDuration(task.OutFile)
	if err != nil {
		log.Println("Failed to get video duration:", err)
		// 可以使用默认有效期，比如 1 小时
		videoDuration = 3600
	}
	expiry := time.Duration(videoDuration) * time.Second
	videoURL, err := k.minio.PresignedGetObject(
		context.Background(),
		"videos",
		objectKey,
		expiry,
		nil,
	)
	if err != nil {
		log.Println("Failed to generate presigned URL:", err)
		return err
	}

	f.Close()
	if err := os.RemoveAll(task.OutFile); err != nil {
		log.Println("Failed to remove file", err)
	}
	if err := os.RemoveAll(task.TmepFile); err != nil {
		log.Println("Failed to remove file", err)
	}
	// 清理 Redis
	idxChunkKey := fmt.Sprintf("%s%s", constants.UPLOAD_CHUNK_IDX, task.UploadID)
	if err := k.redisOp.Del(context.Background(), idxChunkKey); err != nil {
		log.Println("Failed to delete Redis key:", err)
	}

	log.Println("debug")

	// TODO 将上传好的 url 传入到 redis 中，再通过 ws 的初始化传递给前端
	var message = res.VideoStateInit{
		RoomID:    task.RoomID,
		Progress:  0,
		IsPlaying: false,
		Speed:     1.0,
		TimeStamp: time.Now().UnixMilli(),
		VideoUrl:  videoURL.String(),
	}

	stateKey := fmt.Sprintf("%s%s", constants.VIDEO_STATE_ROOM_INIT, task.RoomID)
	messageByte, err := json.Marshal(message)
	if err != nil {
		log.Println("Failed to resolve message", err)
		return err
	}

	ok, err := k.redisOp.SetNx(context.TODO(), stateKey, messageByte, expiry)

	if !ok {
		log.Println("Failed to set videoUrl to redis", err)
	}

	return err
}
