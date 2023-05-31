package api

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/grpc/pb"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/location"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
)

// GrpcSrv server for grpc
var GrpcSrv *grpc.Server

// internal service
var (
	grpcLocationHandler *location.Handler
)

func InjectDependency(locationHub *location.Hub, consumerGroup sarama.ConsumerGroup) {
	// internal dependency
	grpcLocationHandler = location.NewHandler(locationHub, consumerGroup)
}

func InitGRPCServer() {
	srv := grpc.NewServer()
	pb.RegisterLocationServer(srv, grpcLocationHandler)

	GrpcSrv = srv
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		logger.Logger.Fatal("failed to listen tcp", zap.Error(err))
	}

	logger.Logger.Info("starting server", zap.String("port", lis.Addr().String()))
	if err = GrpcSrv.Serve(lis); err != nil {
		logger.Logger.Fatal("failed to server grpc server", zap.Error(err))
	}
}
