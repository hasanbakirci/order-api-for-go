package rabbitmqclient

import (
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/internal/config"
	"github.com/streadway/amqp"
)

type Client struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMqClient(settings config.RabbitMQSettings) (*Client, error) {
	url := settings.Url
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	//defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	//defer ch.Close()
	return &Client{
		Connection: conn,
		Channel:    ch,
	}, err
	return nil, err
}
