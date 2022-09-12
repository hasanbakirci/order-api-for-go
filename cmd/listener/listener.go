package listener

import (
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/internal/config"
	"github.com/hasanbakirci/order-api-for-go/internal/logger"
	"github.com/hasanbakirci/order-api-for-go/internal/order"
	"github.com/hasanbakirci/order-api-for-go/internal/queue"
	"github.com/hasanbakirci/order-api-for-go/pkg/mongoHelper"
	"github.com/hasanbakirci/order-api-for-go/pkg/rabbitmqclient"
	"github.com/hasanbakirci/order-api-for-go/pkg/redisClient"
)

type listener struct {
	rabbitClient   rabbitmqclient.Client
	redisClient    *redisClient.RedisClient
	deleteConsumer queue.DeleteConsumer
	loggerConsumer queue.LoggerConsumer
	redisConsumer  queue.RedisConsumer
}

func NewListener(settings config.Configuration) listener {
	// Mongo
	db, err := mongoHelper.ConnectDb(settings.MongoSettings)
	if err != nil {
		fmt.Println("Db connection error")
	}
	// RabbitMQ
	client := rabbitmqclient.NewRabbitClient(settings)
	client.DeclareExchangeQueueBindings()

	// Redis
	redis := redisClient.NewRedisClient("localhost:6379")

	redisConsumer := queue.NewRedisConsumer(redis)

	repository := order.NewRepository(db)
	service := order.NewService(repository, client)
	deleteConsumer := queue.NewDeleteConsumer(service, client)

	logRepository := logger.NewRepository(db)
	logService := logger.NewLogService(logRepository)
	logConsumer := queue.NewLoggerConsumer(logService, client)

	return listener{
		rabbitClient:   *client,
		redisClient:    redis,
		deleteConsumer: deleteConsumer,
		loggerConsumer: logConsumer,
		redisConsumer:  redisConsumer,
	}
}
func (l listener) Start() {
	go l.loggerConsumer.Consume()
	go l.deleteConsumer.Consume()
	go l.redisConsumer.Consume("redisLog")

}
