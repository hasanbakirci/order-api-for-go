package listener

import (
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/internal/config"
	"github.com/hasanbakirci/order-api-for-go/internal/consumer"
	"github.com/hasanbakirci/order-api-for-go/internal/order"
	"github.com/hasanbakirci/order-api-for-go/pkg/mongoHelper"
	"github.com/hasanbakirci/order-api-for-go/pkg/rabbitmqclient"
)

type listener struct {
	rabbitClient   rabbitmqclient.Client
	deleteConsumer consumer.DeleteConsumer
}

func NewListener(settings config.Configuration) listener {

	db, err := mongoHelper.ConnectDb(settings.MongoSettings)
	if err != nil {
		fmt.Println("Db connection error")
	}
	client, cErr := rabbitmqclient.NewRabbitMqClient(settings.RabbitMQSettings)
	if cErr != nil {
		fmt.Println("RabbitMQ connection error")
	}

	repository := order.NewRepository(db)
	service := order.NewService(repository)
	consumer := consumer.NewDeleteConsumer(service, client)

	return listener{
		rabbitClient:   *client,
		deleteConsumer: consumer,
	}
}
func (l listener) Start() {
	l.deleteConsumer.Consume()
}
