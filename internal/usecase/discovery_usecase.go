package usecase

import (
	"context"
	"math"

	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/repository"
)

type (
	DiscoveryUsecase interface {
		SearchProduct(ctx context.Context, payload dto.SearchProductPayload) (*dto.SearchProductResponse, error)
	}
	discoveryUsecase struct {
		pr repository.ProductRepository
		rr repository.ReviewRepository
	}
)

// SearchProduct implements DiscoveryUsecase.
func (uc *discoveryUsecase) SearchProduct(ctx context.Context, payload dto.SearchProductPayload) (*dto.SearchProductResponse, error) {
	res := &dto.SearchProductResponse{
		Page:       payload.Page,
		SearchTerm: payload.SearchTerm,
	}

	products, err := uc.pr.FindProductBySearchTerm(ctx, payload)
	if err != nil {
		return nil, err
	}

	for i, product := range products {
		rate, err := uc.rr.RateOfProduct(ctx, product.ProductCode)
		if err != nil {
			return nil, err
		}
		rating := rate.RateSum / rate.RateCount
		if math.IsNaN(rating) {
			rating = 0
		}
		products[i].Rating = rating
	}

	res.Products = products
	res.TotalProduct = len(products)

	countPayload := dto.CountProductBySearchTermPayload{
		SearchTerm:  payload.SearchTerm,
		DistrictIDs: payload.DistrictIDs,
		CategoryID:  payload.CategoryID,
		MinPrice:    payload.MinPrice,
		MaxPrice:    payload.MaxPrice,
	}
	count, err := uc.pr.CountProductBySearchTerm(ctx, countPayload)
	if err != nil {
		return nil, err
	}

	res.TotalPage = int(math.Ceil(float64(*count) / 30.0))

	return res, nil
}

func NewDiscoveryUsecase(pr repository.ProductRepository, rr repository.ReviewRepository) DiscoveryUsecase {
	return &discoveryUsecase{
		pr: pr,
		rr: rr,
	}
}
