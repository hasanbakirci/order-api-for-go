package order

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(context.Context, *Order) (string, error)
	Update(ctx context.Context, order Order) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetAll(ctx context.Context) ([]Order, error)
	GetById(ctx context.Context, id string) (*Order, error)
	GetByCustomerId(ctx context.Context, id string) ([]Order, error)
	ChangeStatus(ctx context.Context, id string, status string) (bool, error)
}

type mongoRepository struct {
	collection *mongo.Collection
}

func (m mongoRepository) GetByCustomerId(ctx context.Context, id string) ([]Order, error) {
	cursor, err := m.collection.Find(ctx, bson.M{"customer_id": id})
	orders := make([]Order, 0)
	if err = cursor.All(ctx, &orders); err != nil {
		log.Fatal(err)
	}
	return orders, err
}

func (m mongoRepository) ChangeStatus(ctx context.Context, id string, status string) (bool, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"status": &status,
	}}
	_, err := m.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m mongoRepository) Delete(ctx context.Context, id string) (bool, error) {
	_, err := m.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m mongoRepository) Update(ctx context.Context, order Order) (bool, error) {
	filter := bson.M{"_id": order.Id}
	update := bson.M{"$set": bson.M{
		"customer_id": order.CustomerId,
		"quantity":    order.Quantity,
		"status":      order.Status,
	}}
	_, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m mongoRepository) GetById(ctx context.Context, id string) (order *Order, err error) {
	err = m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	return
}

func (m mongoRepository) GetAll(ctx context.Context) ([]Order, error) {
	cursor, err := m.collection.Find(ctx, bson.M{})
	orders := make([]Order, 0)
	if err = cursor.All(ctx, &orders); err != nil {
		log.Fatal(err)
	}
	return orders, err
}

func (m mongoRepository) Create(ctx context.Context, order *Order) (string, error) {
	order.Id = uuid.New().String()
	order.Status = "Waiting"
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
