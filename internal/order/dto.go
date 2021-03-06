package order

import (
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// swagger:model CreateOrderRequest
type CreateOrderRequest struct {
	CustomerId string               `json:"customer_id" validate:"required"`
	Quantity   int                  `json:"quantity" validate:"required,gt=0,numeric"`
	Price      float32              `bson:"price" validate:"required,gt=0,numeric"`
	Address    CreateAddressRequest `bson:"address" validate:"required"`
	Product    CreateProductRequest `bson:"product" validate:"required"`
}

// swagger:model CreateAddressRequest
type CreateAddressRequest struct {
	AddressLine string `bson:"address_line" validate:"required"`
	City        string `bson:"city" validate:"required"`
	Country     string `bson:"country" validate:"required"`
	CityCode    int    `bson:"city_code" validate:"required,gt=0,numeric"`
}

// swagger:model CreateProductRequest
type CreateProductRequest struct {
	Id       string `bson:"id" validate:"required"`
	ImageUrl string `bson:"image_url" validate:"required"`
	Name     string `bson:"name" validate:"required"`
}

// swagger:model UpdateOrderRequest
type UpdateOrderRequest struct {
	CustomerId string `json:"customer_id" validate:"required"`
	Quantity   int    `json:"quantity" validate:"required,gt=0,numeric"`
	Status     string `json:"status" validate:"required"`
}

// swagger:model ChangeStatusRequest
type ChangeStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

// swagger:model OrderResponse
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

// swagger:model AddressResponse
type AddressResponse struct {
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	Country     string `json:"country"`
	CityCode    int    `json:"city_code"`
}

// swagger:model ProductResponse
type ProductResponse struct {
	Id       string `json:"id"`
	ImageUrl string `json:"image_url"`
	Name     string `json:"name"`
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
	id, _ := uuid.FromBytes(o.Id.Data)
	cid, _ := uuid.FromBytes(o.CustomerId.Data)
	pid, _ := uuid.FromBytes(o.Product.Id.Data)
	return &OrderResponse{
		Id:         id.String(),
		CustomerId: cid.String(),
		Quantity:   o.Quantity,
		Price:      o.Price,
		Status:     o.Status,
		Address:    AddressResponse(o.Address),
		Product: ProductResponse{
			Id:       pid.String(),
			ImageUrl: o.Product.ImageUrl,
			Name:     o.Product.Name,
		},
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
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
