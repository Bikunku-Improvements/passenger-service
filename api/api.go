package api

import (
	"fmt"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/grpc/pb"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/location"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

// GrpcSrv server for grpc
var GrpcSrv *grpc.Server

// KafkaReaderLocation external service for kafka reader
var (
	KafkaReaderLocation *kafka.Reader
)

// internal service
var (
	grpcLocationHandler *location.Handler
)

func InjectDependency() {
	// need to close connection dependency
	KafkaReaderLocation = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{os.Getenv("KAFKA_ADDR")},
		Topic:     "location",
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	// internal dependency
	locationRepository := location.NewRepository(KafkaReaderLocation)
	locationUseCase := location.NewUseCase(locationRepository)
	grpcLocationHandler = location.NewHandler(locationUseCase)
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
