package order

import (
	"context"
	"github.com/pkg/errors"
)

// Servcice
type Service interface {
	Create(context.Context, Order) (string, error)
}

type service struct {
	repository Repository
}

//Create Order Method
func (s service) Create(ctx context.Context, order Order) (string, error) {
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
