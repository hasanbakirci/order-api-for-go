package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func RecoverMiddlewareFunc(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			str := recover()
			c.JSON(http.StatusInternalServerError, str)
		}()
		return handlerFunc(c)
	}
}
