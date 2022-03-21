package order

import (
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CreateOrderRequest struct {
	CustomerId string               `json:"customer_id" validate:"required"`
	Quantity   int                  `json:"quantity" validate:"required,gt=0,numeric"`
	Price      float32              `bson:"price" validate:"required,gt=0,numeric"`
	Address    CreateAddressRequest `bson:"address" validate:"required"`
	Product    CreateProductRequest `bson:"product" validate:"required"`
}
type CreateAddressRequest struct {
	AddressLine string `bson:"address_line" validate:"required"`
	City        string `bson:"city" validate:"required"`
	Country     string `bson:"country" validate:"required"`
	CityCode    int    `bson:"city_code" validate:"required,gt=0,numeric"`
}
type CreateProductRequest struct {
	Id       string `bson:"id" validate:"required"`
	ImageUrl string `bson:"image_url" validate:"required"`
	Name     string `bson:"name" validate:"required"`
}
type UpdateOrderRequest struct {
	CustomerId string `json:"customer_id" validate:"required"`
	Quantity   int    `json:"quantity" validate:"required,gt=0,numeric"`
	Status     string `json:"status" validate:"required"`
}
type ChangeStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type OrderResponse struct {
	Id         primitive.Binary `json:"id"`
	CustomerId primitive.Binary `json:"customer_id"`
	Quantity   int              `json:"quantity"`
	Price      float32          `json:"price"`
	Status     string           `json:"status"`
	Address    AddressResponse  `json:"address"`
	Product    ProductResponse  `json:"product"`
	CreatedAt  time.Time        `json:"createdAt"`
	UpdatedAt  time.Time        `json:"updatedAt"`
}
type AddressResponse struct {
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	Country     string `json:"country"`
	CityCode    int    `json:"city_code"`
}
type ProductResponse struct {
	Id       primitive.Binary `json:"id"`
	ImageUrl string           `json:"image_url"`
	Name     string           `json:"name"`
}

func (c *CreateOrderRequest) ToOrder() *Order {
	cid, _ := uuid.FromString(c.CustomerId)
	pid, _ := uuid.FromString(c.Product.Id)
	return &Order{
		CustomerId: primitive.Binary{3, cid.Bytes()},
		Quantity:   c.Quantity,
		Price:      c.Price,
		Address:    Address(c.Address),
		Product: Product{
			Id:       primitive.Binary{3, pid.Bytes()},
			ImageUrl: c.Product.ImageUrl,
			Name:     c.Product.Name,
		},
	}
}

func (o Order) ToOrderResponse() *OrderResponse {
	return &OrderResponse{
		Id:         o.Id,
		CustomerId: o.CustomerId,
		Quantity:   o.Quantity,
		Price:      o.Price,
		Status:     o.Status,
		Address:    AddressResponse(o.Address),
		Product:    ProductResponse(o.Product),
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
	}
}

func (u UpdateOrderRequest) ToOrder() *Order {
	cid, _ := uuid.FromString(u.CustomerId)
	return &Order{
		CustomerId: primitive.Binary{3, cid.Bytes()},
		Quantity:   u.Quantity,
		Status:     u.Status,
	}
}
