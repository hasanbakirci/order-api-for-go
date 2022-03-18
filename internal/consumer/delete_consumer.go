package consumer

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hasanbakirci/order-api-for-go/internal/order"
	"github.com/hasanbakirci/order-api-for-go/pkg/rabbitmqclient"
	//rabbit "github.com/streadway/amqp"
)

type DeleteConsumer struct {
	service order.Service
	client  *rabbitmqclient.Client
}

func NewDeleteConsumer(s order.Service, c *rabbitmqclient.Client) DeleteConsumer {
	return DeleteConsumer{
		service: s,
		client:  c,
	}
}

func (d DeleteConsumer) Consume() {
	messages, err := d.client.Channel.Consume(
		"Deleted-Customer-Queue",
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
		deletedId, _ := uuid.Parse(string(m.Body))
		ok, _ := d.service.DeleteCustomersOrder(context.Background(), deletedId)
		fmt.Printf("Recived Message: %s || status : %t \n", m.Body, ok)
	}
}
