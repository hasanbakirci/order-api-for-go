package logger

import (
	"context"
	"encoding/json"
	"github.com/hasanbakirci/order-api-for-go/internal/order"
	"github.com/pkg/errors"
)

type Service interface {
	Create(ctx context.Context, message []byte) (bool, error)
}

type service struct {
	logRepo Repository
}

func (s service) Create(ctx context.Context, message []byte) (bool, error) {
	var order order.Order
	e := json.Unmarshal(message, &order)
	if e != nil {
		return false, errors.New("message is nil")
	}
	ok, err := s.logRepo.Create(ctx, order)
	if ok {
		return ok, err
	}
	return ok, err
}

func NewLogService(r Repository) Service {
	return &service{
		logRepo: r,
	}
}
