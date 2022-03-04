package order

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Resource struct {
	service Service
}

func NewResource(s Service) Resource {
	return Resource{service: s}
}

func RegisterHandlers(instance *echo.Echo, api Resource) {
	baseUrl := "order"
	instance.GET("/", func(c echo.Context) error {
		c.JSON(http.StatusOK, "orders")
		return nil
	})
	instance.POST(fmt.Sprintf("%s", baseUrl), api.createOrder)
}

func (r *Resource) createOrder(c echo.Context) error {
	request := new(CreateOrderRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	result, err := r.service.Create(c.Request().Context(), *request.ToOrder())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
