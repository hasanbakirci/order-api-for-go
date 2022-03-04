package order

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(context.Context, *Order) (string, error)
}

type mongoRepository struct {
	collection *mongo.Collection
}

func (m mongoRepository) Create(ctx context.Context, order *Order) (string, error) {
	_, err := m.collection.InsertOne(ctx, order)
	return order.Id, err
}

//------- infra-------------> <- factory -><-- repo>>
//mongodriver -> mongoclient -> database -> collection
func NewRepository(db *mongo.Database) Repository {
	col := db.Collection("orders")
	return &mongoRepository{
		collection: col,
	}
}
