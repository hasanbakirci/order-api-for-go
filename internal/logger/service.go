package logger

import (
	"context"
	"github.com/pkg/errors"
)

type Service interface {
	Create(ctx context.Context, message string) (bool, error)
}

type service struct {
	logRepo Repository
}

func (s service) Create(ctx context.Context, message string) (bool, error) {
	if message == "" {
		return false, errors.New("message is nil")
	}
	ok, err := s.logRepo.Create(ctx, message)
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
