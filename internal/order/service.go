package order

import (
	"context"
	"github.com/pkg/errors"
)

// Servcice
type Service interface {
	Create(context.Context, CreateOrderRequest) (string, error)
	Update(ctx context.Context, request UpdateOrderRequest) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetAll(ctx context.Context) ([]OrderResponse, error)
	GetById(ctx context.Context, id string) (*OrderResponse, error)
	GetByCustomerId(ctx context.Context, id string) ([]OrderResponse, error)
	ChangeStatus(ctx context.Context, request ChangeStatusRequest) (bool, error)
}

type service struct {
	repository Repository
}

func (s service) GetByCustomerId(ctx context.Context, id string) ([]OrderResponse, error) {
	orders, err := s.repository.GetByCustomerId(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "Service customerid error")
	}
	orderResponse := make([]OrderResponse, 0)
	for i := 0; i < len(orders); i++ {
		t := orders[i].ToOrderResponse()
		orderResponse = append(orderResponse, *t)
	}
	return orderResponse, nil

}

func (s service) ChangeStatus(ctx context.Context, request ChangeStatusRequest) (bool, error) {
	result, err := s.repository.ChangeStatus(ctx, request.Id, request.Status)
	return result, err
}

func (s service) Delete(ctx context.Context, id string) (status bool, err error) {
	status, err = s.repository.Delete(ctx, id)
	if err != nil {
		err = errors.Wrap(err, "Delete error")
	}
	return status, err
}

func (s service) Update(ctx context.Context, request UpdateOrderRequest) (status bool, err error) {
	order := request.ToOrder()
	status, err = s.repository.Update(ctx, order)
	return
}

func (s service) GetById(ctx context.Context, id string) (response *OrderResponse, err error) {
	order, err := s.repository.GetById(ctx, id)
	if err != nil {
		err = errors.Wrapf(err, "Service get by id error id: %s", id)
	}
	response = order.ToOrderResponse()
	return

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
