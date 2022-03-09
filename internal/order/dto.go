package order

import "time"

type CreateOrderRequest struct {
	CustomerId string               `json:"customer_id"`
	Quantity   int                  `json:"quantity"`
	Price      float32              `bson:"price"`
	Address    CreateAddressRequest `bson:"address"`
	Product    CreateProductRequest `bson:"product"`
}
type CreateAddressRequest struct {
	AddressLine string `bson:"address_line"`
	City        string `bson:"city"`
	Country     string `bson:"country"`
	CityCode    int    `bson:"city_code"`
}
type CreateProductRequest struct {
	Id       string `bson:"id"`
	ImageUrl string `bson:"image_url"`
	Name     string `bson:"name"`
}

func (c *CreateOrderRequest) ToOrder() *Order {
	return &Order{
		CustomerId: c.CustomerId,
		Quantity:   c.Quantity,
		Price:      c.Price,
		Address:    Address(c.Address),
		Product:    Product(c.Product),
	}
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

type UpdateOrderRequest struct {
	Id         string `json:"id"`
	CustomerId string `json:"customer_id"`
	Quantity   int    `json:"quantity"`
	Status     string `json:"status"`
}

func (u UpdateOrderRequest) ToOrder() *Order {
	return &Order{
		Id:         u.Id,
		CustomerId: u.CustomerId,
		Quantity:   u.Quantity,
		Status:     u.Status,
	}
}

type ChangeStatusRequest struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
