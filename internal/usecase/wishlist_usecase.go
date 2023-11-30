package usecase

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/shopspring/decimal"
)

type (
	WishlistUseCase interface {
		AddProductToWishlist(ctx context.Context, payload *dto.WishlistPayload) error
		DeleteProductFromWishlist(ctx context.Context, payload *dto.WishlistPayload) error
		GetAllProductOfUser(ctx context.Context, userID int64, params *dto.WishlistParams) (*dto.WishlistUserResponse, error)
	}
	wishlistUseCase struct {
		wr repository.WishlistRepository
		pr repository.ProductRepository
		rr repository.ReviewRepository
	}
)

func (wuc *wishlistUseCase) AddProductToWishlist(ctx context.Context, payload *dto.WishlistPayload) error {
	product, err := wuc.pr.FirstProductByCode(ctx, payload.ProductCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.ErrProductNotFound
		}
		return err
	}
	wlist, err := wuc.wr.FirstByUserAndProduct(ctx, payload.UserID, product.ID)
	if err != nil {
		return err
	}
	if wlist != nil {
		return shared.ErrWishlistAlreadyExist
	}
	wishlist := &model.Wishlist{
		AccountID: payload.UserID,
		ProductID: product.ID,
	}
	if err := wuc.wr.Create(ctx, wishlist); err != nil {
		return err
	}
	return nil
}

func (wuc *wishlistUseCase) DeleteProductFromWishlist(ctx context.Context, payload *dto.WishlistPayload) error {
	product, err := wuc.pr.FirstProductByCode(ctx, payload.ProductCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.ErrProductNotFound
		}
		return err
	}

	item, err := wuc.wr.FirstByUserAndProduct(ctx, payload.UserID, product.ID)
	if item == nil {
		return shared.ErrWishlistNotFound
	}
	if item.AccountID != payload.UserID {
		return shared.ErrUnauthorizedUser
	}
	if err != nil {
		return err
	}
	if err = wuc.wr.Delete(ctx, item.ID); err != nil {
		return err
	}
	return nil
}

func (wuc *wishlistUseCase) GetAllProductOfUser(ctx context.Context, userID int64, params *dto.WishlistParams) (*dto.WishlistUserResponse, error) {
	wishlist, err := wuc.wr.FindByAccountID(ctx, userID, params)
	if err != nil {
		return nil, err
	}
	items := make([]dto.WishlistUserResponseItems, 0)
	for _, wl := range wishlist {
		rate, err := wuc.rr.RateOfProduct(ctx, wl.ProductCode)
		if err != nil {
			return nil, err
		}
		rating := rate.RateSum / rate.RateCount
		if math.IsNaN(rating) {
			rating = 0
		}
		items = append(items, dto.WishlistUserResponseItems{
			ID:            wl.ID,
			ProductCode:   wl.ProductCode,
			ProductName:   wl.ProductName,
			ThumbnailURL:  wl.ThumbnailURL,
			BasePrice:     wl.BasePrice.InexactFloat64(),
			Discount:      wl.Discount,
			DiscountPrice: wl.BasePrice.Mul(decimal.NewFromFloat(float64(100-wl.Discount) / 100)).InexactFloat64(),
			ShopName:      wl.ShopName,
			DistrictName:  wl.DistrictName,
			Rating:        rating,
		})
	}
	metadata, err := wuc.wr.FindByAccountIDMetadata(ctx, userID)
	if err != nil {
		return nil, err
	}
	totalData := len(metadata)
	totalPage := math.Ceil(float64(totalData) / constant.WishlistDefaultItems)
	res := &dto.WishlistUserResponse{
		Items:       items,
		CurrentPage: params.Page,
		TotalPage:   int(totalPage),
		TotalData:   totalData,
	}
	return res, nil
}

func NewWishlistUsecase(wr repository.WishlistRepository, pr repository.ProductRepository, rr repository.ReviewRepository) WishlistUseCase {
	return &wishlistUseCase{
		wr: wr,
		pr: pr,
		rr: rr,
	}
}
