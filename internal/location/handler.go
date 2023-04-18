package location

import (
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/grpc/pb"
	"github.com/google/uuid"
)

type Handler struct {
	pb.UnimplementedLocationServer
	broadcaster *Broadcaster
}

func (h *Handler) SubscribeLocation(req *pb.SubscribeLocationRequest, stream pb.Location_SubscribeLocationServer) error {
	requestID := uuid.NewString()

	h.broadcaster.Subscribe(stream, requestID)
	<-stream.Context().Done()
	h.broadcaster.Unsubscribe(requestID)

	return nil
}

func NewHandler(broadcaster *Broadcaster) *Handler {
	return &Handler{
		broadcaster: broadcaster,
	}
}
