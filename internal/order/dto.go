package order

import (
	"github.com/google/uuid"
	"time"
)

type CreateOrderRequest struct {
	CustomerId uuid.UUID            `json:"customer_id" validate:"required"`
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
	Id       uuid.UUID `bson:"id" validate:"required"`
	ImageUrl string    `bson:"image_url" validate:"required"`
	Name     string    `bson:"name" validate:"required"`
}
type UpdateOrderRequest struct {
	//Id         uuid.UUID `json:"id" validate:"required"`
	CustomerId uuid.UUID `json:"customer_id" validate:"required"`
	Quantity   int       `json:"quantity" validate:"required,gt=0,numeric"`
	Status     string    `json:"status" validate:"required"`
}
type ChangeStatusRequest struct {
	//Id     uuid.UUID `json:"id"`
	Status string `json:"status" validate:"required"`
}

type OrderResponse struct {
	Id         string          `json:"id"`
	CustomerId string          `json:"customer_id"`
	Quantity   int             `json:"quantity"`
	Price      float32         `json:"price"`
	Status     string          `json:"status"`
	Address    AddressResponse `json:"address"`
	Product    ProductResponse `json:"product"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}
type AddressResponse struct {
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	Country     string `json:"country"`
	CityCode    int    `json:"city_code"`
}
type ProductResponse struct {
	Id       string `json:"id"`
	ImageUrl string `json:"image_url"`
	Name     string `json:"name"`
}

func (c *CreateOrderRequest) ToOrder() *Order {
	return &Order{
		CustomerId: c.CustomerId.String(),
		Quantity:   c.Quantity,
		Price:      c.Price,
		Address:    Address(c.Address),
		Product: Product{
			Id:       c.Product.Id.String(),
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
	return &Order{
		//Id:         u.Id.String(),
		CustomerId: u.CustomerId.String(),
		Quantity:   u.Quantity,
		Status:     u.Status,
	}
}
