package order

import (
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository interface {
	Create(context.Context, *Order) (primitive.Binary, error)
	Update(ctx context.Context, id primitive.Binary, order Order) (bool, error)
	Delete(ctx context.Context, id primitive.Binary) (bool, error)
	GetAll(ctx context.Context) ([]Order, error)
	GetById(ctx context.Context, id primitive.Binary) (*Order, error)
	GetByCustomerId(ctx context.Context, id primitive.Binary) ([]Order, error)
	ChangeStatus(ctx context.Context, id primitive.Binary, status string) (bool, error)
	DeleteCustomersOrder(ctx context.Context, id primitive.Binary) (bool, error)
}

type mongoRepository struct {
	collection *mongo.Collection
}

func (m mongoRepository) DeleteCustomersOrder(ctx context.Context, id primitive.Binary) (bool, error) {
	filter := bson.M{"customer_id": id}
	deleteResult, err := m.collection.DeleteMany(ctx, filter)
	if deleteResult.DeletedCount > 0 {
		fmt.Printf("deleted orders for customerid:%s,count:%d", id, deleteResult.DeletedCount)
		return true, nil
	}
	return false, err
}

func (m mongoRepository) GetByCustomerId(ctx context.Context, id primitive.Binary) ([]Order, error) {
	cursor, err := m.collection.Find(ctx, bson.M{"customer_id": id})
	orders := make([]Order, 0)
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (m mongoRepository) ChangeStatus(ctx context.Context, id primitive.Binary, status string) (bool, error) {

	filter := bson.M{"_id": id}
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

func (m mongoRepository) Delete(ctx context.Context, id primitive.Binary) (bool, error) {
	deleteResult, err := m.collection.DeleteOne(ctx, bson.M{"_id": id})
	if deleteResult.DeletedCount < 1 {
		return false, err
	}
	return true, nil
}

func (m mongoRepository) Update(ctx context.Context, id primitive.Binary, order Order) (bool, error) {
	filter := bson.M{"_id": id}
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

func (m mongoRepository) GetById(ctx context.Context, id primitive.Binary) (*Order, error) {
	order := new(Order)
	if err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order); err != nil {
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

func (m mongoRepository) Create(ctx context.Context, order *Order) (primitive.Binary, error) {
	order.Id = primitive.Binary{
		Subtype: 3,
		Data:    uuid.NewV3(uuid.NewV4(), "order").Bytes(),
	}
	order.Status = "Waiting"
	order.CreatedAt = time.Now()
	result, err := m.collection.InsertOne(ctx, order)
	if result.InsertedID == nil {
		return primitive.Binary{}, err
	}
	return order.Id, err
}

func NewRepository(db *mongo.Database) Repository {
	col := db.Collection("orders")
	return &mongoRepository{
		collection: col,
	}
}
