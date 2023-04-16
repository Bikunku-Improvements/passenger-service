package location

import (
	"context"
	"fmt"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/reader"
	"github.com/segmentio/kafka-go"
)

type Repository struct {
	reader *kafka.Reader
}

func (r Repository) GetLocation(ctx context.Context, locChan chan Location, errChan chan error) {
	reader := reader.NewKafkaReaderLocation()
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to read message: %v", err)
		}

		var loc Location
		if err = loc.UnmarshalBinary(m.Value); err != nil {
			errChan <- fmt.Errorf("failed to unmarshal location: %v", err)
		}

		locChan <- loc
	}
}

func NewRepository() *Repository {
	return &Repository{}
}
