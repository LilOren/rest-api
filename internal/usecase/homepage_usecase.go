package usecase

import (
	"context"

	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/shopspring/decimal"
)

type (
	HomepageUsecase interface {
		GetRecommendedProducts(ctx context.Context) ([]dto.HomePageProductResponseBody, error)
		GetCartForHome(ctx context.Context, accountId int64) ([]dto.CartHomeResponse, error)
		GetTopCategories(ctx context.Context) ([]dto.HomePageCategoryResponseBody, error)
	}
	homepageUsecase struct {
		pr   repository.ProductRepository
		cr   repository.CartRepository
		rr   repository.ReviewRepository
		catr repository.CategoryRepository
		ccr  repository.CacheRepository
	}
)

// GetTopCategories implements HomepageUsecase.
func (uc *homepageUsecase) GetTopCategories(ctx context.Context) ([]dto.HomePageCategoryResponseBody, error) {
	cats, err := uc.catr.FindByPairCategoryID(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]dto.HomePageCategoryResponseBody, 0)
	for _, c := range cats {
		temp := dto.HomePageCategoryResponseBody{
			TopLevelCategoryID:   c.FirstLevelID,
			ChildLevelCategoryID: c.SecondLevelID,
			CategoryName:         c.FirstLevelName,
			ImageURL:             c.FirstCategoryImageURL,
		}

		res = append(res, temp)
	}

	return res, nil
}

func (uc *homepageUsecase) GetRecommendedProducts(ctx context.Context) ([]dto.HomePageProductResponseBody, error) {
	res, err := uc.ccr.GetRecommendedProduct(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (uc *homepageUsecase) GetCartForHome(ctx context.Context, accountId int64) ([]dto.CartHomeResponse, error) {
	carts, err := uc.cr.FindCartForHome(ctx, accountId)
	if err != nil {
		return nil, err
	}
	cartsRes := make([]dto.CartHomeResponse, 0)
	for _, cart := range carts {
		cartsRes = append(cartsRes, dto.CartHomeResponse{
			ProductName:  cart.ProductName,
			ThumbnailUrl: cart.ThumbnailUrl,
			Price:        cart.BasePrice.Mul(decimal.NewFromFloat((100 - cart.Discount) / 100)).InexactFloat64(),
			Quantity:     cart.Quantity,
		})
	}
	return cartsRes, nil
}

func NewHomepageUsecase(
	pr repository.ProductRepository,
	cr repository.CartRepository,
	rr repository.ReviewRepository,
	catr repository.CategoryRepository,
	ccr repository.CacheRepository,
) HomepageUsecase {
	return &homepageUsecase{
		pr:   pr,
		cr:   cr,
		rr:   rr,
		catr: catr,
		ccr:  ccr,
	}
}
