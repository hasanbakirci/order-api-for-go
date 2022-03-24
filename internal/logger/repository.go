package logger

import (
	"context"
	"github.com/hasanbakirci/order-api-for-go/internal/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, message order.Order) (bool, error)
}
type loggerRepository struct {
	collection *mongo.Collection
}

func (l loggerRepository) Create(ctx context.Context, message order.Order) (bool, error) {
	result, err := l.collection.InsertOne(ctx, bson.M{"message": message}) //
	if result.InsertedID == nil {
		return false, err
	}
	return true, err
}

func NewRepository(db *mongo.Database) Repository {
	col := db.Collection("loggers")
	return &loggerRepository{
		collection: col,
	}
}
