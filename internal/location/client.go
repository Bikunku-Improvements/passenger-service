package location

import (
	"encoding/json"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/grpc/pb"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Client struct {
	hub *Hub

	// The websocket connection.
	conn pb.Location_SubscribeLocationServer

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *Client) writePump() {
	for v := range c.send {
		var loc Event
		err := json.Unmarshal(v, &loc)
		if err != nil {
			logger.Logger.Error("failed to marshal data", zap.Error(err))
		}

		logger.Logger.Info("message receive", zap.Duration("latency", time.Now().Sub(loc.CreatedAt)))
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
				return
			default:
				logger.Logger.Error("error from grpc", zap.Error(err))
			}
		}
	}
}
