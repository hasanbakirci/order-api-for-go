package order

import (
	"github.com/hasanbakirci/order-api-for-go/pkg/response"
	"github.com/hasanbakirci/order-api-for-go/pkg/validationHelper"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Controller struct {
	service Service
}

func NewController(s Service) Controller {
	return Controller{service: s}
}

func RegisterHandlers(instance *echo.Echo, api Controller) {
	instance.GET("/", func(c echo.Context) error {
		c.JSON(http.StatusOK, "orders")
		return nil
	})
	instance.POST("api/order", api.createOrder)
	instance.GET("api/order", api.getallOrders)
	instance.GET("api/order/:id", api.getbyid)
	instance.PUT("api/order/:id", api.updateOrder)
	instance.DELETE("api/order/:id", api.deleteOrder)
	instance.GET("api/order/customer/:id", api.getbyCustomer)
	instance.PUT("api/order/status/:id", api.changeStatus)
}

func (r *Controller) createOrder(c echo.Context) error {
	request := new(CreateOrderRequest)
	//if err := c.Bind(request); err != nil {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}
	if _, err := validationHelper.Validate(c, request); err != nil {
		//return c.JSON(http.StatusBadRequest, err)
		return response.ErrorResponse(c, 400, 4002, err.Error())
	}
	result, err := r.service.Create(c.Request().Context(), *request)
	if err != nil {
		//return c.JSON(http.StatusInternalServerError, err.Error())
		return response.ErrorResponse(c, 400, 4000, err.Error())
	}
	//return c.JSON(http.StatusCreated, result)
	return response.SuccessResponse(c, 201, result, "success")
}

func (r *Controller) getallOrders(c echo.Context) error {
	orders, err := r.service.GetAll(c.Request().Context())
	if err != nil {
		//return c.JSON(http.StatusNotFound, "")
		return response.ErrorResponse(c, 401, 4011, err.Error())
	}
	//return c.JSON(http.StatusOK, orders)
	return response.SuccessResponse(c, 200, orders, "success")
}

func (r *Controller) getbyid(c echo.Context) error {
	//id, _ := uuid.Parse(c.Param("id"))
	id, _ := uuid.FromString(c.Param("id"))
	order, err := r.service.GetById(c.Request().Context(), primitive.Binary{Subtype: 3, Data: id.Bytes()})
	if err != nil {
		//return c.JSON(http.StatusNotFound, err.Error())
		return response.ErrorResponse(c, 400, 4001, err.Error())
	}
	return response.SuccessResponse(c, 200, order, "success")
}

func (r *Controller) updateOrder(c echo.Context) error {
	id, _ := uuid.FromString(c.Param("id"))
	request := new(UpdateOrderRequest)
	//if err := c.Bind(request); err != nil {
	//	c.JSON(http.StatusBadRequest, err.Error())
	//}
	if _, err := validationHelper.Validate(c, request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	ok, err := r.service.Update(c.Request().Context(), primitive.Binary{Subtype: 3, Data: id.Bytes()}, *request)
	if !ok {
		//c.JSON(http.StatusInternalServerError, err.Error())
		return response.ErrorResponse(c, 400, 4002, err.Error())
	}
	//return c.JSON(http.StatusOK, ok)
	return response.SuccessResponse(c, 200, ok, "success")
}

func (r *Controller) deleteOrder(c echo.Context) error {
	id, _ := uuid.FromString(c.Param("id"))
	ok, err := r.service.Delete(c.Request().Context(), primitive.Binary{Subtype: 3, Data: id.Bytes()})
	if !ok {
		//return c.JSON(http.StatusNotFound, err.Error())
		return response.ErrorResponse(c, 401, 4012, err.Error())
	}
	//return c.JSON(http.StatusAccepted, ok)
	return response.SuccessResponse(c, 200, ok, "success")
}

func (r *Controller) getbyCustomer(c echo.Context) error {
	id, _ := uuid.FromString(c.Param("id"))
	orders, err := r.service.GetByCustomerId(c.Request().Context(), primitive.Binary{Subtype: 3, Data: id.Bytes()})
	if err != nil {
		//return c.JSON(http.StatusNotFound, err.Error())
		return response.ErrorResponse(c, 401, 4013, err.Error())
	}
	//return c.JSON(http.StatusOK, orders)
	return response.SuccessResponse(c, 200, orders, "success")
}

func (r *Controller) changeStatus(c echo.Context) error {
	id, _ := uuid.FromString(c.Param("id"))
	request := new(ChangeStatusRequest)
	//if err := c.Bind(request); err != nil {
	//	c.JSON(http.StatusBadRequest, err.Error())
	//}
	if _, err := validationHelper.Validate(c, request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	ok, err := r.service.ChangeStatus(c.Request().Context(), primitive.Binary{Subtype: 3, Data: id.Bytes()}, *request)
	if !ok {
		//return c.JSON(http.StatusNotFound, err.Error())
		return response.ErrorResponse(c, 401, 4014, err.Error())
	}
	//return c.JSON(http.StatusOK, ok)
	return response.SuccessResponse(c, 200, ok, "success")
}
