package consumer

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"sync"
)

func NewSarama(brokerID []string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	brokers := brokerID
	groupID := "location-consumer-group-1"

	log.Println("Starting a new Sarama consumer")

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Printf("NewConsumerGroupFromClient(%s) error: %v", groupID, err)
	}

	return client, nil
}

func Consume(ctx context.Context, client sarama.ConsumerGroup, consumer *Handler) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, []string{"location"}, consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")
}
