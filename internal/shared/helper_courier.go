package shared

import "github.com/lil-oren/rest/internal/dto"

func ServiceChanger(ro dto.RajaOngkirGetCostHTTPQueries) dto.RajaOngkirGetCostHTTPQueries {
	ro.CourierService = "CTC"
	return ro
}
