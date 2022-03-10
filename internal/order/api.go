package order

import (
	"fmt"
	"github.com/hasanbakirci/order-api-for-go/pkg/validationHelper"
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
	baseUrl := "api/order"
	instance.GET("/", func(c echo.Context) error {
		c.JSON(http.StatusOK, "orders")
		return nil
	})
	instance.POST(fmt.Sprintf("%s", baseUrl), api.createOrder)
	instance.GET(fmt.Sprintf(baseUrl), api.getallOrders)
	instance.GET("api/order/:id", api.getbyid)
	instance.PUT("api/order", api.updateOrder)
	instance.DELETE("api/order/:id", api.deleteOrder)
	instance.GET("api/order/customer/:id", api.getbyCustomer)
	instance.PUT("api/order/status", api.changeStatus)
}

func (r *Resource) createOrder(c echo.Context) error {
	request := new(CreateOrderRequest)
	//if err := c.Bind(request); err != nil {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}
	if _, err := validationHelper.Validate(c, request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	result, err := r.service.Create(c.Request().Context(), *request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, result)
}

func (r *Resource) getallOrders(c echo.Context) error {
	orders, err := r.service.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusNotFound, "")
	}
	return c.JSON(http.StatusOK, orders)
}

func (r *Resource) getbyid(c echo.Context) error {
	id := c.Param("id")
	order, err := r.service.GetById(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, order)
}

func (r *Resource) updateOrder(c echo.Context) error {
	request := new(UpdateOrderRequest)
	//if err := c.Bind(request); err != nil {
	//	c.JSON(http.StatusBadRequest, err.Error())
	//}
	if _, err := validationHelper.Validate(c, request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	result, err := r.service.Update(c.Request().Context(), *request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func (r *Resource) deleteOrder(c echo.Context) error {
	id := c.Param("id")
	status, err := r.service.Delete(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusAccepted, status)
}

func (r *Resource) getbyCustomer(c echo.Context) error {
	id := c.Param("id")
	orders, err := r.service.GetByCustomerId(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, orders)
}

func (r *Resource) changeStatus(c echo.Context) error {
	request := new(ChangeStatusRequest)
	//if err := c.Bind(request); err != nil {
	//	c.JSON(http.StatusBadRequest, err.Error())
	//}
	if _, err := validationHelper.Validate(c, request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	result, err := r.service.ChangeStatus(c.Request().Context(), *request)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
