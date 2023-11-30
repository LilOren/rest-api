package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
)

type (
	RajaOngkirRepository interface {
		GetCost(ctx context.Context, query dto.RajaOngkirGetCostHTTPQueries) (*float64, error)
	}
	rajaOngkirRepository struct {
		config dependency.Config
	}
)

// GetCost implements RajaOngkirRepository.
func (r *rajaOngkirRepository) GetCost(ctx context.Context, query dto.RajaOngkirGetCostHTTPQueries) (*float64, error) {
	formData := url.Values{}
	formData.Set("origin", strconv.Itoa(query.OriginCityID))
	formData.Set("destination", strconv.Itoa(query.DestinationCityID))
	formData.Set("weight", strconv.Itoa(query.Weight))
	formData.Set("courier", query.CourierCode)

	roURL := fmt.Sprintf("%s%s", r.config.ThirdParty.RajaOngkirBaseURL, "/cost")
	req, err := http.NewRequest(http.MethodPost, roURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add(constant.RajaOngkirAPIKeyHeader, r.config.ThirdParty.RajaOngkirAPIKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	result := new(dto.RajaOngkirGetCostHTTPResponse)
	if err := json.NewDecoder(res.Body).Decode(result); err != nil {
		return nil, err
	}

	firstResult := result.RajaOngkir.Results[0]

	cost := float64(0)
	for _, c := range firstResult.Costs {
		if c.Service == query.CourierService {
			cost = float64(c.Cost[0].Value)
		}
	}

	return &cost, nil
}

func NewRajaOngkirRepository(config dependency.Config) RajaOngkirRepository {
	return &rajaOngkirRepository{
		config: config,
	}
}
