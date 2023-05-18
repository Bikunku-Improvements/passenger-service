package main

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/api"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/consumer"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/location"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	hub := location.NewHub()
	go hub.Run()

	// start consumer
	broker := strings.Split(os.Getenv("KAFKA_ADDR"), ";")
	consumerGroup, err := consumer.NewSarama(broker)
	if err != nil {
		log.Panicf("failed to init sarama consumer: %v", err)
	}

	go consumer.Consume(context.Background(), consumerGroup, consumer.Handler{Hub: hub, Ready: make(chan bool)})

	// start grpc server
	go func() {
		api.InjectDependency(hub, consumerGroup)
		api.InitGRPCServer()
		api.StartGRPCServer()
	}()

	<-ctx.Done()
	stop()

	defer func() {
		log.Printf("closing consumer...")
		consumerGroup.Close()
	}()

	defer func() {
		log.Printf("closing grpc server...")
		api.GrpcSrv.GracefulStop()
	}()
}
