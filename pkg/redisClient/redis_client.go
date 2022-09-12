package redisClient

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	redisClient *redis.Client
}

func NewRedisClient(host string) *RedisClient {
	client := redis.NewClient(&redis.Options{Addr: host})
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
	return &RedisClient{redisClient: client}
}

func (redis RedisClient) Publish(channel string, message interface{}) {
	err := redis.redisClient.Publish(context.Background(), channel, message).Err()
	if err != nil {
		panic(err)
	}
}

func (redis RedisClient) Subscribe(channel string) *redis.PubSub {
	ctx := context.Background()
	subs := redis.redisClient.Subscribe(ctx, channel)
	//for {
	//	msg, err := subs.ReceiveMessage(ctx)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(msg.Channel, msg.Payload)
	//}
	return subs
}

func (redis RedisClient) Set(key string, value interface{}) {
	v, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	//redis.redisClient.LPush(context.TODO(), key+"00", v)
	redis.redisClient.Set(context.TODO(), key, v, 100*time.Second)
}
