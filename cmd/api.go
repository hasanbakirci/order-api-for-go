/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/internal/config"
	"github.com/hasanbakirci/order-api-for-go/internal/order"
	"github.com/hasanbakirci/order-api-for-go/pkg/echoExtensions"
	"github.com/hasanbakirci/order-api-for-go/pkg/middleware"
	"github.com/hasanbakirci/order-api-for-go/pkg/mongoHelper"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"time"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "My Restfull Order API",
	Long:  `This project is a simple order api`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("api called")
	//},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	var port string
	var cfgFile string
	apiCmd.PersistentFlags().StringVarP(&port, "port", "p", "1994", "Restfull Service Port")
	apiCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.dev", "config file (default is $HOME/.golang-api.yaml)")

	ApiConfig, err := config.GetAllValues("./config/", cfgFile)
	if err != nil {
		panic(err)
	}

	apiCmd.Run = func(cmd *cobra.Command, args []string) {
		//application bootstrapper
		instance := echo.New()
		//custom middleware
		instance.Use(middleware.RecoverMiddlewareFunc)

		db, err := mongoHelper.ConnectDb(ApiConfig.MongoSettings)
		if err != nil {
			fmt.Println("Db connection error")
		}

		//Register handlers -> resource -> service -> repository -> mongodb
		repository := order.NewRepository(db)
		service := order.NewService(repository)
		controller := order.NewController(service)
		order.RegisterHandlers(instance, controller)

		fmt.Println("Api starting")
		if err := instance.Start(fmt.Sprintf(":%s", port)); err != nil {
			fmt.Println("Api fatal error")
		}
		echoExtensions.Shutdown(instance, time.Second*2)
	}

}
