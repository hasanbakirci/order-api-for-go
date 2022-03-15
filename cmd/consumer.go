/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/internal/config"
	"github.com/hasanbakirci/order-api-for-go/internal/order"
	"github.com/hasanbakirci/order-api-for-go/pkg/mongoHelper"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

type Consumer struct {
	service order.Service
}

func NewConsumer(s order.Service) Consumer {
	return Consumer{service: s}
}

//func (c *Consumer) DeleteCustomers(id string) (result bool,err error){
//	result, err = c.service.DeleteCustomersOrder(context.TODO(),id)
//	return
//}

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "This project is consume the deleted customer",

	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("consumer called")
	//},
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

		repository := order.NewRepository(db)
		service := order.NewService(repository)
		consumer := NewConsumer(service)

		fmt.Println("Consumer Application")
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer conn.Close()
		fmt.Println("Succesfuly Connected to our RabbitMQ instance")

		ch, err := conn.Channel()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer ch.Close()

		msgs, err := ch.Consume(
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
				fmt.Printf("Recived Message: %s\n", d.Body)
				consumer.service.DeleteCustomersOrder(context.TODO(), string(d.Body))
			}
		}()

		fmt.Println("Succesfully connected to our RabbitMq instance")
		fmt.Println("[*] - waiting go messages")
		<-forever
	}

}
