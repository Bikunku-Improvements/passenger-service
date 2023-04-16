package reader

import (
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"os"
	"strings"
)

func NewKafkaReaderLocation() *kafka.Reader {
	addr := strings.Split(os.Getenv("KAFKA_ADDR"), ";")
	KafkaReaderLocation := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     addr,
		Topic:       "location",
		StartOffset: kafka.LastOffset,
		GroupID:     "location-consumer-group-" + uuid.NewString(),
	})

	return KafkaReaderLocation
}
