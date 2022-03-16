package order

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository interface {
	Create(context.Context, *Order) (string, error)
	Update(ctx context.Context, id uuid.UUID, order Order) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	GetAll(ctx context.Context) ([]Order, error)
	GetById(ctx context.Context, id uuid.UUID) (*Order, error)
	GetByCustomerId(ctx context.Context, id uuid.UUID) ([]Order, error)
	ChangeStatus(ctx context.Context, id uuid.UUID, status string) (bool, error)
	DeleteCustomersOrder(ctx context.Context, id uuid.UUID) (bool, error)
}

type mongoRepository struct {
	collection *mongo.Collection
}

func (m mongoRepository) DeleteCustomersOrder(ctx context.Context, id uuid.UUID) (bool, error) {
	filter := bson.M{"customer_id": id.String()}
	deleteResult, err := m.collection.DeleteMany(ctx, filter)
	if deleteResult.DeletedCount > 0 {
		fmt.Printf("deleted orders for customerid:%s,count:%d", id, deleteResult.DeletedCount)
		return true, nil
	}
	return false, err
}

func (m mongoRepository) GetByCustomerId(ctx context.Context, id uuid.UUID) ([]Order, error) {
	cursor, err := m.collection.Find(ctx, bson.M{"customer_id": id.String()})
	orders := make([]Order, 0)
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (m mongoRepository) ChangeStatus(ctx context.Context, id uuid.UUID, status string) (bool, error) {

	filter := bson.M{"_id": id.String()}
	update := bson.M{"$set": bson.M{
		"status":    &status,
		"updatedAt": time.Now(),
	}}
	updateResult, err := m.collection.UpdateOne(ctx, filter, update)
	if updateResult.ModifiedCount < 1 {
		return false, err
	}
	return true, nil
}

func (m mongoRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	deleteResult, err := m.collection.DeleteOne(ctx, bson.M{"_id": id.String()})
	if deleteResult.DeletedCount < 1 {
		return false, err
	}
	return true, nil
}

func (m mongoRepository) Update(ctx context.Context, id uuid.UUID, order Order) (bool, error) {
	filter := bson.M{"_id": id.String()}
	update := bson.M{"$set": bson.M{
		"customer_id": order.CustomerId,
		"quantity":    order.Quantity,
		"status":      order.Status,
		"updatedAt":   time.Now(),
	}}
	updateResult, err := m.collection.UpdateOne(ctx, filter, update)
	if updateResult.ModifiedCount < 1 {
		return false, err
	}
	return true, nil
}

func (m mongoRepository) GetById(ctx context.Context, id uuid.UUID) (*Order, error) {
	order := new(Order)
	if err := m.collection.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&order); err != nil {
		return nil, err
	}
	return order, nil
}

func (m mongoRepository) GetAll(ctx context.Context) ([]Order, error) {
	cursor, err := m.collection.Find(ctx, bson.M{})
	orders := make([]Order, 0)
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (m mongoRepository) Create(ctx context.Context, order *Order) (string, error) {
	order.Id = uuid.New().String()
	order.Status = "Waiting"
	order.CreatedAt = time.Now()
	result, err := m.collection.InsertOne(ctx, order)
	if result.InsertedID == nil {
		return "", err
	}
	return order.Id, err
}

func NewRepository(db *mongo.Database) Repository {
	col := db.Collection("orders")
	return &mongoRepository{
		collection: col,
	}
}
