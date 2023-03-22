package location

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type useCase interface {
	GetLocation(ctx context.Context) (*Location, error)
}

type Handler struct {
	pb.UnimplementedLocationServer
	useCase useCase
}

func (h Handler) GetLocation(req *pb.GetLocationRequest, stream pb.Location_GetLocationServer) error {
	for {
		location, err := h.useCase.GetLocation(stream.Context())
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

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

func NewHandler(useCase useCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}
