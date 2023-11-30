package dto

type (
	RajaOngkirGetCostHTTPQueries struct {
		OriginCityID      int
		DestinationCityID int
		Weight            int
		CourierCode       string
		CourierService    string
	}
	RajaOngkirGetCostHTTPResponse struct {
		RajaOngkir struct {
			Results []struct {
				Costs []struct {
					Service string `json:"service"`
					Cost    []struct {
						Value int `json:"value"`
					} `json:"cost"`
				} `json:"costs"`
			} `json:"results"`
		} `json:"rajaongkir"`
	}
)
