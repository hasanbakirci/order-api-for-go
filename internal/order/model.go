package order

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type (
	Order struct {
		Id         primitive.Binary `bson:"_id"`
		CustomerId primitive.Binary `bson:"customer_id"`
		Quantity   int              `bson:"quantity"`
		Price      float32          `bson:"price"`
		Status     string           `bson:"status"`
		Address    Address          `bson:"address"`
		Product    Product          `bson:"product"`
		CreatedAt  time.Time        `bson:"createdAt"`
		UpdatedAt  time.Time        `bson:"updatedAt"`
	}
	Address struct {
		AddressLine string `bson:"address_line"`
		City        string `bson:"city"`
		Country     string `bson:"country"`
		CityCode    int    `bson:"city_code"`
	}
	Product struct {
		Id       primitive.Binary `bson:"id"`
		ImageUrl string           `bson:"image_url"`
		Name     string           `bson:"name"`
	}
)
