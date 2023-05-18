package api

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/grpc/pb"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/location"
	"google.golang.org/grpc"
	"log"
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
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := GrpcSrv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
