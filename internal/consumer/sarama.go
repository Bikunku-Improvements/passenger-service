package consumer

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"sync"
)

type ConsumerGroup interface {
	sarama.ConsumerGroupHandler
	SetReady(chan bool)
	GetReady() chan bool
}

func NewSarama(brokerID []string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	brokers := brokerID
	groupID := "location-consumer-group-1"

	log.Printf("Starting a new Sarama consumer with broker: %v\n", brokers)

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
				log.Panicf("Error from consumer: %v", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")
}
