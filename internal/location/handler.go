package location

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type useCase interface {
	GetLocation(ctx context.Context, locChan chan Location, errChan chan error)
}

type Handler struct {
	pb.UnimplementedLocationServer
	useCase useCase
}

func (h *Handler) GetLocation(req *pb.GetLocationRequest, stream pb.Location_GetLocationServer) error {
	locChan := make(chan Location, 1)
	errChan := make(chan error, 1)

	go h.useCase.GetLocation(stream.Context(), locChan, errChan)

	for {
		if location, ok := <-locChan; ok {
			log.Printf("message received with latency %s", time.Now().Sub(location.CreatedAt))
			if err := stream.Send(&pb.GetLocationResponse{
				Long:      location.Long,
				Lat:       location.Lat,
				CreatedAt: location.CreatedAt.Format(time.RFC3339),
				BusId:     location.BusID,
			}); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		}
	}
}

func NewHandler(useCase useCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}
