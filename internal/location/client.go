package location

import (
	"encoding/json"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type Client struct {
	hub *Hub

	// The websocket connection.
	conn pb.Location_SubscribeLocationServer

	// Buffered channel of outbound messages.
	send chan []byte
}

//func (c *Client) SetReady(bool chan bool) {
//	c.ready = bool
//}
//
//func (c *Client) GetReady() chan bool {
//	return c.ready
//}
//
//func (c *Client) Setup(_ sarama.ConsumerGroupSession) error {
//	close(c.ready)
//	return nil
//}
//
//func (c *Client) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
//
//func (c *Client) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
//	for {
//		select {
//		case message := <-claim.Messages():
//			session.MarkMessage(message, "")
//			c.hub.Broadcast <- message.Value
//
//		case <-session.Context().Done():
//			return nil
//		}
//	}
//}

func (c *Client) writePump() {
	for v := range c.send {
		var loc Event
		err := json.Unmarshal(v, &loc)
		if err != nil {
			log.Printf("failed to marshal data: %v\n", err.Error())
		}

		log.Printf("message received with latency: %s", time.Now().Sub(loc.CreatedAt))
		err = c.conn.Send(&pb.SubscribeLocationResponse{
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
				log.Printf("error from grpc: %v\n", err.Error())
			}
		}
	}
}
