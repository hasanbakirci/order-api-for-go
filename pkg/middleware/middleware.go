package middleware

import (
	"github.com/hasanbakirci/order-api-for-go/pkg/response"
	"github.com/labstack/echo/v4"
)

func RecoverMiddlewareFunc(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
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
		return handlerFunc(c)
	}
}
