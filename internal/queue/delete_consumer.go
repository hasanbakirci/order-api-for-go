package queue

import (
	"context"
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/internal/order"
	"github.com/hasanbakirci/order-api-for-go/pkg/rabbitmqclient"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//rabbit "github.com/streadway/amqp"
)

type DeleteConsumer struct {
	service  order.Service
	client   *rabbitmqclient.Client
	producer *Producer
}

func NewDeleteConsumer(s order.Service, c *rabbitmqclient.Client, p *Producer) DeleteConsumer {
	return DeleteConsumer{
		service:  s,
		client:   c,
		producer: p,
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
		deletedId, _ := uuid.FromString(string(m.Body))
		orders, _ := d.service.GetByCustomerId(context.Background(), primitive.Binary{3, deletedId.Bytes()})
		for i := 0; i < len(orders); i++ {
			d.service.Delete(context.Background(), orders[i].Id)
			d.producer.QueueDeclare()
			d.producer.Publish(orders[i])
		}
		fmt.Printf("Recived Message: %s", m.Body)
	}
}
