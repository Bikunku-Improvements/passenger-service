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

type Event struct {
	BusID     uint      `json:"bus_id"`
	Number    int       `json:"number"`
	Plate     string    `json:"plate"`
	Status    string    `json:"status"`
	Route     string    `json:"route"`
	IsActive  bool      `json:"isActive"`
	Long      float64   `json:"long"`
	Lat       float64   `json:"lat"`
	Speed     float64   `json:"speed"`
	Heading   float64   `json:"heading"`
	CreatedAt time.Time `json:"created_at"`
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

		var loc Event
		err := json.Unmarshal(msg.Value, &loc)
		if err != nil {
			return err
		}

		b.sync.Lock()
		for _, stream := range b.clientStream {
			log.Printf("message received with latency: %s", time.Now().Sub(loc.CreatedAt))
			err := stream.Send(&pb.SubscribeLocationResponse{
				BusId:     uint64(loc.BusID),
				Number:    int64(loc.Number),
				Plate:     loc.Plate,
				Status:    loc.Status,
				Route:     loc.Route,
				IsActive:  loc.IsActive,
				Long:      float32(loc.Long),
				Lat:       float32(loc.Lat),
				Speed:     float32(loc.Speed),
				Heading:   float32(loc.Heading),
				CreatedAt: loc.CreatedAt.Format(time.RFC3339Nano),
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
