package clients

import (
	"errors"
	"github.com/valyala/fasthttp"
)

func ValidateCustomer(id string) (bool, error) {
	url := "https://localhost:5001/api/customers/validate/" + id
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := fasthttp.Do(req, resp)
	if err != nil {
		return false, errors.New("client get failed")
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return false, errors.New("customer id not found")
	}
	return true, nil
}
