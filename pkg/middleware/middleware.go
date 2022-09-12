package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/pkg/redisClient"
	"github.com/hasanbakirci/order-api-for-go/pkg/response"
	"github.com/labstack/echo/v4"
)

func RecoverMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			//str := recover()
			//c.JSON(http.StatusInternalServerError, str)
			if r := recover(); r != nil {
				switch t := r.(type) {
				case response.ErrorDetails:
					response.ErrorResponse(c, t.StatusCode, t.ErrorCode, t.Message)
				default:
					response.ErrorResponse(c, 500, 5000, r)
				}
				//c.JSON(http.StatusInternalServerError, err)
				//response.ErrorResponse(c, 500, 5000, r)
			}
		}()
		return next(c)
	}
}

func TokenHandlerMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if len(token) > 0 {
			if token != "Bearer ABC" {
				return response.ErrorResponse(c, 401, 40101, "token error")
			}
			return next(c)
		}
		return response.ErrorResponse(c, 401, 40101, "token error")
	}
}

func LoggerMiddlewareFunc(redis *redisClient.RedisClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			log := fmt.Sprintf("Request method: '%s' Url: '%s' ", c.Request().Method, c.Request().RequestURI)
			if err := next(c); err != nil {
				c.Error(err)
			}
			log += fmt.Sprintf("Response status: '%d'", c.Response().Status)
			fmt.Println("---> ", log)
			body, _ := json.Marshal(log)
			redis.Publish("redisLog", body)
			return next(c)
		}
	}
}

type MyLog struct {
	RequestMethod  string `json:"request_method"`
	RequestUrl     string `json:"request_method"`
	ResponseStatus int    `json:"request_method"`
}
