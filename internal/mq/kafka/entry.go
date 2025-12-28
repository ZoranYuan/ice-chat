package kafka

import (
	"context"
	"ice-chat/config"
	"ice-chat/pkg/ws"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

var sy sync.Once

// KafkaClient Kafka客户端
type KafkaClient struct {
	writer      *kafka.Writer
	reader      *kafka.Reader
	roomManager *ws.RoomManager
}

func NewKafkaClient(roomManager *ws.RoomManager) *KafkaClient {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Conf.Kafka.Brokers...),
		Topic:        config.Conf.Kafka.Topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: time.Duration(config.Conf.Kafka.Producer.WriteTimeout),
		BatchSize:    1,
		RequiredAcks: kafka.RequireOne,
		Async:        true,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Conf.Kafka.Brokers,
		Topic:          config.Conf.Kafka.Topic,
		GroupID:        config.Conf.Kafka.GroupID,
		MinBytes:       config.Conf.Kafka.Consumer.MinBytes,
		MaxBytes:       config.Conf.Kafka.Consumer.MaxBytes,
		ReadBackoffMin: time.Duration(config.Conf.Kafka.Consumer.ReadBackoffMin),
		MaxWait:        1 * time.Millisecond,
	})

	return &KafkaClient{
		writer:      writer,
		reader:      reader,
		roomManager: roomManager,
	}
}

func (k *KafkaClient) Consume() {
	defer func() {
		_ = k.reader.Close()
	}()

	for {
		_, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("kafka fail to consume : %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// TODO 根据不同的 Topic 去找对应的服务
	}
}

func (k *KafkaClient) Produce(ctx context.Context, msg []byte) error {
	return k.writer.WriteMessages(ctx,
		kafka.Message{
			Value: msg,
			Time:  time.Now(),
		},
	)
}
