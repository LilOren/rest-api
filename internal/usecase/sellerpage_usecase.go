package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
)

type (
	SellerPageUsecase interface {
		GetSellerDetail(ctx context.Context, shopName string, payload dto.SearchSellerProductPayload) (*dto.GetSellerDetailResponseBody, error)
	}
	sellerPageUsecase struct {
		spr repository.SellerPageRepository
		rr  repository.ReviewRepository
	}
)

func (spu *sellerPageUsecase) GetSellerDetail(ctx context.Context, shopName string, payload dto.SearchSellerProductPayload) (*dto.GetSellerDetailResponseBody, error) {
	detail, err := spu.spr.FirstSellerDetail(ctx, shopName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.ErrShopNotFound
		}
		return nil, shared.ErrFindShopDetail
	}

	sellerDetail := &dto.GetSellerDetailResponseBody{
		ShopName:     detail.ShopName,
		ProductCount: detail.ProductCount,
		Years:        fmt.Sprintf("%d years", detail.Years),
	}

	if detail.Years == 1 {
		sellerDetail.Years = "1 year"
	}

	if detail.Years < 1 {
		sellerDetail.Years = "< 1 year"
	}

	if payload.Page < 0 {
		return nil, shared.ErrInvalidPage
	}

	res := &dto.SearchSellerProductResponse{
		Page:       payload.Page,
		SearchTerm: payload.SearchTerm,
	}

	categories, err := spu.spr.FindCategoryBySellerId(ctx, detail.SellerId)
	if err != nil {
		return nil, shared.ErrFindCategory
	}

	products, err := spu.spr.FindSellerProductBySearchTerm(ctx, payload, detail.SellerId)
	if err != nil {
		return nil, shared.ErrFindProduct
	}

	count, err := spu.spr.CountProductBySearchTerm(ctx, payload, detail.SellerId)
	if err != nil {
		return nil, shared.ErrCountProduct
	}

	payload.Page = 1
	payload.SortBy = "most_purchased"
	payload.SortDesc = true
	payload.SearchTerm = ""
	payload.CategoryName = ""

	bestSeller, err := spu.spr.FindSellerProductBySearchTerm(ctx, payload, detail.SellerId)
	if err != nil {
		return nil, err
	}

	if len(bestSeller) > 6 {
		bestSeller = bestSeller[:6]
	}

	best := make([]dto.SearchSellerProductResponseItem, 0)

	for _, v := range bestSeller {
		data, err := spu.rr.RateOfProduct(ctx, v.ProductCode)
		if err != nil {
			return nil, err
		}
		if data.RateCount == 0 {
			data.RateCount = 1
		}
		b := dto.SearchSellerProductResponseItem{
			ProductCode:    v.ProductCode,
			ProductName:    v.ProductName,
			ThumbnailURL:   v.ThumbnailURL,
			BasePrice:      v.BasePrice,
			DiscountPrice:  v.DiscountPrice,
			Discount:       v.Discount,
			DistrictName:   v.DistrictName,
			Rating:         data.RateSum / data.RateCount,
			CountPurchased: v.CountPurchased,
		}
		best = append(best, b)
	}

	prod := make([]dto.SearchSellerProductResponseItem, 0)

	for _, v := range products {
		data, err := spu.rr.RateOfProduct(ctx, v.ProductCode)
		if err != nil {
			return nil, err
		}
		if data.RateCount == 0 {
			data.RateCount = 1
		}
		p := dto.SearchSellerProductResponseItem{
			ProductCode:    v.ProductCode,
			ProductName:    v.ProductName,
			ThumbnailURL:   v.ThumbnailURL,
			BasePrice:      v.BasePrice,
			DiscountPrice:  v.DiscountPrice,
			Discount:       v.Discount,
			DistrictName:   v.DistrictName,
			Rating:         data.RateSum / data.RateCount,
			CountPurchased: v.CountPurchased,
		}
		prod = append(prod, p)
	}

	res.TotalProduct = *count
	res.TotalPage = int(math.Ceil(float64(res.TotalProduct) / 20.0))

	sellerDetail.Products = prod
	sellerDetail.Pagination = *res
	sellerDetail.Categories = categories
	sellerDetail.BestSeller = best

	return sellerDetail, nil
}

func NewSellerPageUsecase(spr repository.SellerPageRepository, rr repository.ReviewRepository) SellerPageUsecase {
	return &sellerPageUsecase{
		spr: spr,
		rr:  rr,
	}
}
