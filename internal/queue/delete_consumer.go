package queue

import (
	"context"
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/internal/order"
	"github.com/hasanbakirci/order-api-for-go/pkg/rabbitmqclient"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
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
	messages, err := d.client.CreateChannel().Consume(
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
		deletedId, _ := uuid.FromString(stringFormat(string(m.Body)))
		orders, _ := d.service.GetByCustomerId(context.Background(), primitive.Binary{3, deletedId.Bytes()})
		for i := 0; i < len(orders); i++ {
			oid, _ := uuid.FromString(orders[i].Id)
			d.service.Delete(context.Background(), primitive.Binary{3, oid.Bytes()})
		}
		fmt.Printf("delete consumer -> customer id: %s , order lenght: %d", m.Body, len(orders))
	}
}

func stringFormat(str string) string {
	if strings.Contains(str, `"`) {
		newStr := strings.Split(str, `"`)
		return newStr[1]
	}
	return str
}
