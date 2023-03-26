package location

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Repository struct {
	reader *kafka.Reader
}

func (r Repository) GetLocation(ctx context.Context) (*Location, error) {
	m, err := r.reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %v", err)
	}
	log.Printf("message received from %v topic with value %v", m.Topic, string(m.Value))

	var location Location
	if err = location.UnmarshalBinary(m.Value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal location: %v", err)
	}

	log.Printf("time latency, %s", time.Now().Sub(location.CreatedAt))

	return &location, nil
}

func NewRepository(reader *kafka.Reader) *Repository {
	return &Repository{reader: reader}
}
