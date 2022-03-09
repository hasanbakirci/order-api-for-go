package order

import "time"

type (
	Order struct {
		Id         string    `bson:"_id"`
		CustomerId string    `bson:"customer_id"`
		Quantity   int       `bson:"quantity"`
		Price      float32   `bson:"price"`
		Status     string    `bson:"status"`
		Address    Address   `bson:"address"`
		Product    Product   `bson:"product"`
		CreatedAt  time.Time `bson:"createdAt"`
		UpdatedAt  time.Time `bson:"updatedAt"`
	}
	Address struct {
		AddressLine string `bson:"address_line"`
		City        string `bson:"city"`
		Country     string `bson:"country"`
		CityCode    int    `bson:"city_code"`
	}
	Product struct {
		Id       string `bson:"id"`
		ImageUrl string `bson:"image_url"`
		Name     string `bson:"name"`
	}
)
