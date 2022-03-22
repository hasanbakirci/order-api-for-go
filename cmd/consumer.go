/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/hasanbakirci/order-api-for-go/cmd/listener"
	"github.com/hasanbakirci/order-api-for-go/internal/config"
	"github.com/spf13/cobra"
)

// consumerCmd represents the queue command
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

		forever := make(chan bool)

		l := listener.NewListener(*ApiConfig)
		go l.Start()

		<-forever

	}
}
