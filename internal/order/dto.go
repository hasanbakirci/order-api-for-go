package order

import "github.com/google/uuid"

type CreateOrderRequest struct {
	CustomerId string `json:"customer_id"`
	Quantity   int    `json:"quantity"`
}

func (c *CreateOrderRequest) ToOrder() *Order {
	return &Order{
		Id:         uuid.New().String(),
		CustomerId: c.CustomerId,
		Quantity:   c.Quantity,
	}
}
