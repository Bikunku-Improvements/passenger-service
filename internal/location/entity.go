package location

import (
	"encoding/json"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/grpc/pb"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/consumer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sync"
	"time"
)

type Location struct {
	Long      string    `json:"long"`
	Lat       string    `json:"lat"`
	BusID     string    `json:"bus_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (l *Location) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, l)
	if err != nil {
		return err
	}

	return nil
}

type Broadcaster struct {
	sync         sync.Mutex
	clientStream map[string]pb.Location_SubscribeLocationServer
}

func (b *Broadcaster) Subscribe(stream pb.Location_SubscribeLocationServer, id string) {
	b.sync.Lock()
	b.clientStream[id] = stream
	b.sync.Unlock()
}

func (b *Broadcaster) Unsubscribe(id string) {
	b.sync.Lock()
	if _, ok := b.clientStream[id]; ok {
		delete(b.clientStream, id)
	}
	b.sync.Unlock()
}

func (b *Broadcaster) Broadcast(consumer consumer.Handler) error {
	for {
		msg := <-consumer.Message

		var loc Location
		err := loc.UnmarshalBinary(msg.Value)
		if err != nil {
			return err
		}

		log.Printf("message received with latency: %s", time.Now().Sub(loc.CreatedAt))

		b.sync.Lock()
		for _, stream := range b.clientStream {
			err := stream.Send(&pb.SubscribeLocationResponse{
				Long:      loc.Long,
				Lat:       loc.Lat,
				CreatedAt: loc.CreatedAt.Format(time.RFC3339),
				BusId:     loc.BusID,
			})

			if s, ok := status.FromError(err); ok {
				switch s.Code() {
				case codes.OK:
					// do nothing
				case codes.Unavailable, codes.Canceled, codes.DeadlineExceeded:
					// log.Println("client terminated")
				default:
					return err
				}
			}
		}
		b.sync.Unlock()
	}
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clientStream: make(map[string]pb.Location_SubscribeLocationServer),
	}
}
