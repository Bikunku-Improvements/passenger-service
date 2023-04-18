package consumer

import "github.com/Shopify/sarama"

type Handler struct {
	ready   chan bool
	Message chan *sarama.ConsumerMessage
}

func (c *Handler) Setup(_ sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *Handler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			session.MarkMessage(message, "")
			c.Message <- message
		case <-session.Context().Done():
			return nil
		}
	}
}

func NewHandler() *Handler {
	return &Handler{
		ready:   make(chan bool),
		Message: make(chan *sarama.ConsumerMessage),
	}
}
