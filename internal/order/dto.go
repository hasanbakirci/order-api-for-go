package order

type CreateOrderRequest struct {
	CustomerId string `json:"customer_id"`
	Quantity   int    `json:"quantity"`
}

func (c *CreateOrderRequest) ToOrder() *Order {
	return &Order{
		CustomerId: c.CustomerId,
		Quantity:   c.Quantity,
	}
}

type OrderResponse struct {
	Id         string `json:"id"`
	CustomerId string `json:"customer_id"`
	Quantity   int    `json:"quantity"`
	Status     string `json:"status"`
}

func (o Order) ToOrderResponse() *OrderResponse {
	return &OrderResponse{
		Id:         o.Id,
		CustomerId: o.CustomerId,
		Quantity:   o.Quantity,
		Status:     o.Status,
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
