package order

import (
	"context"
	"github.com/hasanbakirci/order-api-for-go/internal/clients"
	rabbit "github.com/hasanbakirci/order-api-for-go/pkg/rabbitmqclient"
	"github.com/hasanbakirci/order-api-for-go/pkg/response"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Servcice
type Service interface {
	Create(context.Context, CreateOrderRequest) (string, error)
	Update(ctx context.Context, id primitive.Binary, request UpdateOrderRequest) (bool, error)
	Delete(ctx context.Context, id primitive.Binary) (bool, error)
	GetAll(ctx context.Context) ([]OrderResponse, error)
	GetById(ctx context.Context, id primitive.Binary) (*OrderResponse, error)
	GetByCustomerId(ctx context.Context, id primitive.Binary) ([]OrderResponse, error)
	ChangeStatus(ctx context.Context, id primitive.Binary, request ChangeStatusRequest) (bool, error)
	DeleteCustomersOrder(ctx context.Context, id primitive.Binary) (bool, error)
}

type service struct {
	repository   Repository
	rabbitClient *rabbit.Client
}

func (s service) DeleteCustomersOrder(ctx context.Context, id primitive.Binary) (bool, error) {
	result, err := s.repository.DeleteCustomersOrder(ctx, id)
	return result, err
}

func (s service) GetByCustomerId(ctx context.Context, id primitive.Binary) ([]OrderResponse, error) {
	orders, err := s.repository.GetByCustomerId(ctx, id)
	orderResponse := make([]OrderResponse, 0)
	for i := 0; i < len(orders); i++ {
		t := orders[i].ToOrderResponse()
		orderResponse = append(orderResponse, *t)
	}
	if len(orderResponse) > 0 {
		return orderResponse, err
	}
	return nil, errors.New("customer id not found")

}

func (s service) ChangeStatus(ctx context.Context, id primitive.Binary, request ChangeStatusRequest) (bool, error) {
	if _, e := s.GetById(ctx, id); e != nil {
		return false, e
	}
	result, err := s.repository.ChangeStatus(ctx, id, request.Status)
	if err != nil {
		err = errors.Wrap(err, "Service error")
		return false, err
	}
	return result, nil
}

func (s service) Delete(ctx context.Context, id primitive.Binary) (bool, error) {
	order, e := s.GetById(ctx, id)
	if e != nil {
		return false, e
	}
	result, err := s.repository.Delete(ctx, id)
	if result {
		s.rabbitClient.PublishMessage("Log-Order-Exchange", "", *order)
		return result, nil
	}
	return result, err
}

func (s service) Update(ctx context.Context, id primitive.Binary, request UpdateOrderRequest) (bool, error) {
	order := request.ToOrder()
	result, err := s.repository.Update(ctx, id, *order)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s service) GetById(ctx context.Context, id primitive.Binary) (*OrderResponse, error) {
	order, err := s.repository.GetById(ctx, id)
	if err != nil {
		err = errors.Wrap(err, "Service")
		return nil, err
	}
	response := order.ToOrderResponse()
	return response, nil

}

func (s service) GetAll(ctx context.Context) ([]OrderResponse, error) {
	orders, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Service getall error")
	}
	orderResponse := make([]OrderResponse, 0)
	for i := 0; i < len(orders); i++ {
		t := orders[i].ToOrderResponse()
		orderResponse = append(orderResponse, *t)
	}
	return orderResponse, nil
}

//Create Order Method
func (s service) Create(ctx context.Context, request CreateOrderRequest) (string, error) {
	status, _ := clients.ValidateCustomer(request.CustomerId)
	if status == false {
		//return primitive.Binary{}, e
		panic(response.CustomPanic(400, 4005, "customer client error"))
	}
	order := *request.ToOrder()
	id, err := s.repository.Create(ctx, &order)
	if err != nil {
		return "", errors.Wrap(err, "Service: Failed to create order")
		//wrap i?? katmandaki hatay?? kaybetmeden yukar?? ????karabilir.
	}
	result, _ := uuid.FromBytes(id.Data)
	return result.String(), nil
}

//ServiceFactory
func NewService(r Repository, c *rabbit.Client) Service {
	if r == nil {
		return nil
	}
	return &service{
		repository:   r,
		rabbitClient: c,
	}
}
