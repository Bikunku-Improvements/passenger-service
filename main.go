package main

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/api"
	"github.com/joho/godotenv"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	go func() {
		api.InjectDependency()
		api.InitGRPCServer()
		api.StartGRPCServer()
	}()

	<-ctx.Done()
	stop()
	api.GrpcSrv.GracefulStop()
	api.KafkaReaderLocation.Close()
}
