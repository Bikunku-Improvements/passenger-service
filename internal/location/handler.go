package location

import (
	"github.com/Shopify/sarama"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/grpc/pb"
)

type Handler struct {
	pb.UnimplementedLocationServer
	hub *Hub
	cg  sarama.ConsumerGroup
}

func (h *Handler) SubscribeLocation(req *pb.SubscribeLocationRequest, stream pb.Location_SubscribeLocationServer) error {
	client := &Client{hub: h.hub, conn: stream, send: make(chan []byte, 256)}
	client.hub.register <- client

	client.writePump()

	<-stream.Context().Done()
	client.hub.unregister <- client

	return nil
}

func NewHandler(hub *Hub, cg sarama.ConsumerGroup) *Handler {
	return &Handler{
		hub: hub,
		cg:  cg,
	}
}
