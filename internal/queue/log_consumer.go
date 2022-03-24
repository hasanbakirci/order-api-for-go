package queue

import (
	"context"
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/internal/logger"
	"github.com/hasanbakirci/order-api-for-go/pkg/rabbitmqclient"
)

type LoggerConsumer struct {
	service logger.Service
	client  *rabbitmqclient.Client
}

func NewLoggerConsumer(s logger.Service, c *rabbitmqclient.Client) LoggerConsumer {
	return LoggerConsumer{
		service: s,
		client:  c,
	}
}
func (l LoggerConsumer) Consume() {
	messages, err := l.client.CreateChannel().Consume(
		"Log-Order-Queue",
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for m := range messages {
		ok, err := l.service.Create(context.Background(), m.Body)
		fmt.Printf("\nlog consumer -> status: %t || err: %s ", ok, err)
	}

}
