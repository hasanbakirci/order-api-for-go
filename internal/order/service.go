package order

import (
	"context"
	"github.com/google/uuid"
	"github.com/hasanbakirci/order-api-for-go/internal/clients"
	"github.com/pkg/errors"
)

// Servcice
type Service interface {
	Create(context.Context, CreateOrderRequest) (string, error)
	Update(ctx context.Context, id uuid.UUID, request UpdateOrderRequest) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	GetAll(ctx context.Context) ([]OrderResponse, error)
	GetById(ctx context.Context, id uuid.UUID) (*OrderResponse, error)
	GetByCustomerId(ctx context.Context, id uuid.UUID) ([]OrderResponse, error)
	ChangeStatus(ctx context.Context, id uuid.UUID, request ChangeStatusRequest) (bool, error)
	DeleteCustomersOrder(ctx context.Context, id uuid.UUID) (bool, error)
}

type service struct {
	repository Repository
}

func (s service) DeleteCustomersOrder(ctx context.Context, id uuid.UUID) (bool, error) {
	result, err := s.repository.DeleteCustomersOrder(ctx, id)
	return result, err
}

func (s service) GetByCustomerId(ctx context.Context, id uuid.UUID) ([]OrderResponse, error) {
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

func (s service) ChangeStatus(ctx context.Context, id uuid.UUID, request ChangeStatusRequest) (bool, error) {
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

func (s service) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	if _, e := s.GetById(ctx, id); e != nil {
		return false, e
	}
	result, err := s.repository.Delete(ctx, id)
	if result {
		return result, nil
	}
	return result, err
}

func (s service) Update(ctx context.Context, id uuid.UUID, request UpdateOrderRequest) (bool, error) {
	if _, e := s.GetById(ctx, id); e != nil {
		return false, e
	}
	order := request.ToOrder()
	result, err := s.repository.Update(ctx, id, *order)
	if err != nil {
		err = errors.Wrap(err, "Service error")
		return false, err
	}
	return result, nil
}

func (s service) GetById(ctx context.Context, id uuid.UUID) (*OrderResponse, error) {
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
	status, e := clients.ValidateCustomer(request.CustomerId.String())
	if status == false {
		return "", e
	}
	order := *request.ToOrder()
	id, err := s.repository.Create(ctx, &order)
	if err != nil {
		return "", errors.Wrap(err, "Service: Failed to create order")
		//wrap iç katmandaki hatayı kaybetmeden yukarı çıkarabilir.
	}
	return id, nil
}

//ServiceFactory
func NewService(r Repository) Service {
	if r == nil {
		return nil
	}
	return &service{
		repository: r,
	}
}
