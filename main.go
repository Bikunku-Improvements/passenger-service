package main

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/api"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/consumer"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/location"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logger.Logger, _ = config.Build()

	err := godotenv.Load()
	if err != nil {
		logger.Logger.Error("failed to load .env", zap.Error(err))
	}

	hub := location.NewHub()
	go hub.Run()

	// start consumer
	broker := strings.Split(os.Getenv("KAFKA_ADDR"), ";")
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	consumerGroup, err := consumer.NewSarama(broker, username, password)
	if err != nil {
		logger.Logger.Fatal("failed to start consumer", zap.Error(err))
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
		logger.Logger.Info("closing consumer")
		consumerGroup.Close()
	}()

	defer func() {
		logger.Logger.Info("closing grpc")
		api.GrpcSrv.GracefulStop()
	}()

	defer func() {
		logger.Logger.Info("sync logger")
		logger.Logger.Sync()
	}()
}
