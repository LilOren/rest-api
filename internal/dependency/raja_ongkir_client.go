package dependency

import "github.com/go-resty/resty/v2"

func NewRajaOngkirClient() *resty.Client {
	return resty.New()
}
