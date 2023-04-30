package main

import (
	"context"
	"github.com/Shopify/sarama"
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

	// start consumer
	broker := strings.Split(os.Getenv("KAFKA_ADDR"), ";")
	consumerGroup, err := consumer.NewSarama(broker)
	if err != nil {
		log.Panicf("failed to init sarama consumer: %v", err)
	}

	saramaHandler := consumer.NewHandler()
	go func(consumerGroup sarama.ConsumerGroup, handler *consumer.Handler) {
		consumer.Consume(ctx, consumerGroup, handler)
		if err != nil {
			log.Panicf("failed to consumer message: %v", err)
		}
	}(consumerGroup, saramaHandler)

	// start broadcaster
	broadcastHandler := location.NewBroadcaster()
	go func(broadcastHandler *location.Broadcaster, consumerHandler consumer.Handler) {
		err := broadcastHandler.Broadcast(consumerHandler)
		if err != nil {
			log.Panicf("failed to start broadcast: %v", err)
		}
	}(broadcastHandler, *saramaHandler)

	// start grpc server
	go func(broadcaster *location.Broadcaster) {
		api.InjectDependency(broadcaster)
		api.InitGRPCServer()
		api.StartGRPCServer()
	}(broadcastHandler)

	<-ctx.Done()
	stop()

	defer func(consumerGroup sarama.ConsumerGroup, handler *consumer.Handler) {
		log.Printf("closing consumer...")
		err = consumerGroup.Close()
		close(handler.Message)
		if err != nil {
			log.Printf("failed to close consumer: %v", err)
		}
	}(consumerGroup, saramaHandler)

	defer func() {
		log.Printf("closing grpc server...")
		api.GrpcSrv.GracefulStop()
	}()
}
