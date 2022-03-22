package queue

import (
	"encoding/json"
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/pkg/rabbitmqclient"
	"github.com/streadway/amqp"
)

type Producer struct {
	client *rabbitmqclient.Client
}

func NewProducer(c *rabbitmqclient.Client) Producer {
	return Producer{client: c}
}

func (p Producer) ExchangeDeclare(exchangeName string, exchangeType string) {
	err := p.client.Channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("producer: failed to declare exchange")
	}
}
func (p Producer) QueueDeclare() {
	p.client.Channel.QueueDeclare(
		"Order-Log-Queue",
		false,
		false,
		false,
		false,
		nil,
	)
}
func (p Producer) Publish(message interface{}) {
	body, _ := json.Marshal(message)

	err := p.client.Channel.Publish(
		"",
		"Order-Log-Queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		fmt.Println("producer: failed to publish message")
	}
}
