/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hasanbakirci/order-api-for-go/internal/config"
	"github.com/hasanbakirci/order-api-for-go/internal/order"
	"github.com/hasanbakirci/order-api-for-go/pkg/mongoHelper"
	"github.com/hasanbakirci/order-api-for-go/pkg/rabbitmqclient"
	"github.com/spf13/cobra"
)

type Consumer struct {
	service order.Service
}

func NewConsumer(s order.Service) Consumer {
	return Consumer{service: s}
}

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "This project is consume the deleted customer",
}

func init() {
	rootCmd.AddCommand(consumerCmd)

	var cfgFile string
	consumerCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.dev", "config file (default is $HOME/.golang-api.yaml)")

	ApiConfig, err := config.GetAllValues("./config/", cfgFile)
	if err != nil {
		panic(err)
	}

	consumerCmd.Run = func(cmd *cobra.Command, args []string) {

		db, err := mongoHelper.ConnectDb(ApiConfig.MongoSettings)
		if err != nil {
			fmt.Println("Db connection error")
		}
		client, cErr := rabbitmqclient.NewRabbitMqClient(ApiConfig.RabbitMQSettings)
		if cErr != nil {
			fmt.Println("RabbitMQ connection error")
		}

		repository := order.NewRepository(db)
		service := order.NewService(repository)
		consumer := NewConsumer(service)

		defer client.Connection.Close()
		defer client.Channel.Close()

		q, err := client.Channel.QueueDeclare(
			"Deleted-Customer-Queue",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Println(q, " created")

		msgs, err := client.Channel.Consume(
			"Deleted-Customer-Queue",
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		forever := make(chan bool)
		go func() {
			for d := range msgs {
				deletedId, _ := uuid.Parse(string(d.Body))
				ok, _ := consumer.service.DeleteCustomersOrder(context.Background(), deletedId)
				fmt.Printf("Recived Message: %s || status : %t \n", d.Body, ok)
			}
		}()
		<-forever
	}

}
