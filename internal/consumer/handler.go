package consumer

import (
	"github.com/Shopify/sarama"
	"github.com/TA-Aplikasi-Pengiriman-Barang/passenger-service/internal/location"
)

type Handler struct {
	Hub   *location.Hub
	Ready chan bool
}

func (c *Handler) Setup(_ sarama.ConsumerGroupSession) error {
	close(c.Ready)
	return nil
}

func (c *Handler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			session.MarkMessage(message, "")
			c.Hub.Broadcast <- message.Value

		case <-session.Context().Done():
			return nil
		}
	}
}
