package rabbitmqclient

import (
	"encoding/json"
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/internal/config"
	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"
	"time"
)

type Client struct {
	connection *amqp.Connection
	config     config.Configuration
}

func NewRabbitClient(config config.Configuration) *Client {
	c := createConnection(config.RabbitConfig)
	return &Client{
		connection: c,
		config:     config,
	}
}

func (c *Client) CloseConnection() {
	c.connection.Close()
}

func (c *Client) CreateChannel() *amqp.Channel {
	channel, err := c.connection.Channel()
	if err != nil {
		channel.Close()
		log.Panicf("Channel could not created. Terminating. Error details: %s", err.Error())
	}
	return channel
}

func (c *Client) PublishMessage(exchangeName string, routingKey string, message interface{}) {
	body, _ := json.Marshal(message)
	err := c.CreateChannel().Publish(
		exchangeName,
		routingKey,
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

func (c *Client) DeclareExchangeQueueBindings() {
	channel := c.CreateChannel()
	declareExchange(channel, c.config.QueuesConfig.Order.LogCreated)
	declareQueue(channel, c.config.QueuesConfig.Order.LogCreated)
	bindQueue(channel, c.config.QueuesConfig.Order.LogCreated)

	declareExchange(channel, c.config.QueuesConfig.Order.OrderDeleted)
	declareQueue(channel, c.config.QueuesConfig.Order.OrderDeleted)
	bindQueue(channel, c.config.QueuesConfig.Order.OrderDeleted)
}

func getConnectionUrl(config config.RabbitConfig) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.VirtualHost)
}

func createConnection(rabbitConfig config.RabbitConfig) *amqp.Connection {
	amqpConfig := amqp.Config{
		Properties: amqp.Table{
			"connection_name": rabbitConfig.ConnectionName,
		},
		Heartbeat: 30 * time.Second,
	}
	connectionUrl := getConnectionUrl(rabbitConfig)
	connection, err := amqp.DialConfig(connectionUrl, amqpConfig)
	if err != nil {
		_ = connection.Close()
		log.Panicf("Client cannogt deserialize. Terminating. Error details: %s", err.Error())
	}
	log.Printf("RabbitMQ connected. Host: %s, Port: %d, Virtual Host: %s", rabbitConfig.Host, rabbitConfig.Port, rabbitConfig.VirtualHost)
	return connection
}

func declareExchange(channel *amqp.Channel, queueConfig config.QueueConfig) {
	err := channel.ExchangeDeclare(queueConfig.Exchange, queueConfig.ExchangeType, true, false, false, false, nil)
	if err != nil {
		log.Panicf("Exchange could not declared. Terminating. Error details: %s", err.Error())
	}
}

func declareQueue(channel *amqp.Channel, queueConfig config.QueueConfig) {
	_, err := channel.QueueDeclare(queueConfig.Queue, true, false, false, false, nil)
	if err != nil {
		log.Panicf("Queue could not declared. Terminating. Error details: %s", err.Error())
	}
}

func bindQueue(channel *amqp.Channel, queueConfig config.QueueConfig) {
	err := channel.QueueBind(queueConfig.Queue, queueConfig.RoutingKey, queueConfig.Exchange, false, nil)
	if err != nil {
		log.Panicf("Binding could not defined. Terminating. Error details: %s", err.Error())
	}
}
