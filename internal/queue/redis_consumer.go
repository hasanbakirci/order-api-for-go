package queue

import (
	"context"
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/pkg/redisClient"
)

type RedisConsumer struct {
	redisClient *redisClient.RedisClient
}

func NewRedisConsumer(redis *redisClient.RedisClient) RedisConsumer {
	return RedisConsumer{redisClient: redis}
}

func (redis RedisConsumer) Consume(channel string) {
	subs := redis.redisClient.Subscribe(channel)
	for {
		msg, err := subs.ReceiveMessage(context.Background())
		if err != nil {
			panic(err)
		}
		redis.redisClient.Set("order-api-redis", msg.Payload)
		fmt.Println(msg.Channel, msg.Payload)
	}
}
