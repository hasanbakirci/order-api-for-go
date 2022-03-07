package order

type (
	Order struct {
		Id         string `bson:"_id"`
		CustomerId string `json:"customer_id"`
		Quantity   int    `json:"quantity"`
		//Price      float32   `json:"price"`
		Status string `json:"status"`
		//Address    Address   `json:"address"`
		//Product    Product   `json:"product"`
		//CreatedAt  time.Time `json:"CreatedAt"`
		//UpdatedAt  time.Time `json:"UpdatedAt"`
	}
	Address struct {
		AddressLine string `json:"address_line"`
		City        string `json:"city"`
		Country     string `json:"country"`
		CityCode    int    `json:"city_code"`
	}
	Product struct {
		Id       string `json:"id"`
		ImageUrl string `json:"image_url"`
		Name     string `json:"name"`
	}
)
