package listener

import (
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/internal/config"
	"github.com/hasanbakirci/order-api-for-go/internal/logger"
	"github.com/hasanbakirci/order-api-for-go/internal/order"
	"github.com/hasanbakirci/order-api-for-go/internal/queue"
	"github.com/hasanbakirci/order-api-for-go/pkg/mongoHelper"
	"github.com/hasanbakirci/order-api-for-go/pkg/rabbitmqclient"
)

type listener struct {
	rabbitClient   rabbitmqclient.Client
	deleteConsumer queue.DeleteConsumer
	loggerConsumer queue.LoggerConsumer
}

func NewListener(settings config.Configuration) listener {

	db, err := mongoHelper.ConnectDb(settings.MongoSettings)
	if err != nil {
		fmt.Println("Db connection error")
	}
	client := rabbitmqclient.NewRabbitClient(settings)
	client.DeclareExchangeQueueBindings()

	repository := order.NewRepository(db)
	service := order.NewService(repository, client)
	deleteConsumer := queue.NewDeleteConsumer(service, client)

	logRepository := logger.NewRepository(db)
	logService := logger.NewLogService(logRepository)
	logConsumer := queue.NewLoggerConsumer(logService, client)

	return listener{
		rabbitClient:   *client,
		deleteConsumer: deleteConsumer,
		loggerConsumer: logConsumer,
	}
}
func (l listener) Start() {
	go l.loggerConsumer.Consume()
	go l.deleteConsumer.Consume()

}
