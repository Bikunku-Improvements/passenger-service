package consumer

import (
	"context"
	"crypto/tls"
	"github.com/Shopify/sarama"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/logger"
	"go.uber.org/zap"
	"strings"
	"sync"
)

type ConsumerGroup interface {
	sarama.ConsumerGroupHandler
	SetReady(chan bool)
	GetReady() chan bool
}

func NewSarama(brokerID []string, username, password string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	config.Net.SASL.Enable = true
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.Net.SASL.Handshake = true
	config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256

	tlsConfig := tls.Config{}
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tlsConfig

	brokers := brokerID
	groupID := "location-consumer-group-1"

	logger.Logger.Info("starting consumer", zap.String("brokers", strings.Join(brokers, ";")))

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Consume(ctx context.Context, client sarama.ConsumerGroup, consumer Handler) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, []string{"location"}, &consumer); err != nil {
				logger.Logger.Fatal("error when consume message", zap.Error(err))
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	logger.Logger.Info("consumer up and running")
}
