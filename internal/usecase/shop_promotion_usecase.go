package usecase

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/shopspring/decimal"
)

type (
	PromotionUsecase interface {
		GetAllPromoFromShop(ctx context.Context, shopID int64, params *dto.PromotionByShopParams) (*dto.PromotionByShopResponse, error)
		AddShopPromotion(ctx context.Context, payload *dto.UpsertShopPromotionPayload) error
		UpdateShopPromotion(ctx context.Context, payload *dto.UpsertShopPromotionPayload, promoID int64) error
		DuplicateShopPromotion(ctx context.Context, promoID int64, sellerID int64) error
		DeleteShopPromotion(ctx context.Context, promoID int64, sellerID int64) error
	}
	promotionUsecase struct {
		pr repository.PromotionRepository
		sr repository.ShopRepository
	}
)

func (uc *promotionUsecase) GetAllPromoFromShop(ctx context.Context, shopID int64, params *dto.PromotionByShopParams) (*dto.PromotionByShopResponse, error) {
	promos, err := uc.pr.FindByShopID(ctx, shopID, params)
	if err != nil {
		return nil, err
	}
	res := make([]dto.PromotionByShopResponseItems, 0)
	for _, promo := range promos {
		res = append(res, dto.PromotionByShopResponseItems{
			ID:           promo.ID,
			Name:         promo.Name,
			ExactPrice:   promo.ExactPrice.Float64,
			Percentage:   promo.Percentage.Float64,
			MinimumSpend: promo.MinimumSpend.InexactFloat64(),
			Quota:        promo.Quota,
			StartedAt:    promo.StartedAt.Format("2006-01-02 15:04:05"),
			ExpiredAt:    promo.ExpiredAt.Format("2006-01-02 15:04:05"),
		})
	}
	metadata, err := uc.pr.FindByShopIDMetadata(ctx, shopID, params)
	if err != nil {
		return nil, err
	}
	totalData := len(metadata)
	totalPage := math.Ceil(float64(totalData) / constant.PromotionShopDefaultItems)
	resp := &dto.PromotionByShopResponse{
		Items:       res,
		TotalData:   totalData,
		TotalPage:   int(totalPage),
		CurrentPage: params.Page,
	}
	return resp, nil
}

func (uc *promotionUsecase) AddShopPromotion(ctx context.Context, payload *dto.UpsertShopPromotionPayload) error {
	promo, err := makePromotionModel(ctx, payload, uc.sr)
	if err != nil {
		return err
	}
	return uc.pr.Create(ctx, promo)
}

func (uc *promotionUsecase) UpdateShopPromotion(ctx context.Context, payload *dto.UpsertShopPromotionPayload, promoID int64) error {
	shop, err := uc.sr.FirstShopById(ctx, int(payload.SellerID))
	if err != nil {
		return err
	}
	promo, err := uc.pr.FirstByID(ctx, promoID)
	if err != nil {
		return err
	}
	if promo == nil {
		return shared.ErrPromoNotFound
	}
	if promo.ShopID != shop.ID {
		return shared.ErrUnauthorizedUser
	}
	promoModel, err := makePromotionModel(ctx, payload, uc.sr)
	if err != nil {
		return err
	}
	promoModel.ID = promoID
	return uc.pr.Update(ctx, promoModel)
}

func makePromotionModel(ctx context.Context, payload *dto.UpsertShopPromotionPayload, sr repository.ShopRepository) (*model.Promotion, error) {
	shop, err := sr.FirstShopById(ctx, int(payload.SellerID))
	if err != nil {
		return nil, err
	}

	startedAt, err := time.ParseInLocation("2006-01-02 15:04:05", payload.StartedAt, time.UTC)
	if err != nil {
		return nil, err
	}
	expiredAt, err := time.ParseInLocation("2006-01-02 15:04:05", payload.ExpiredAt, time.UTC)
	if err != nil {
		return nil, err
	}
	if expiredAt.Before(startedAt) {
		return nil, shared.ErrExpiredBeforeStart
	}

	exactPrice := sql.NullFloat64{}
	if err := exactPrice.Scan(payload.ExactPrice); err != nil {
		return nil, err
	}
	percentage := sql.NullFloat64{}
	if err := percentage.Scan(payload.Percentage); err != nil {
		return nil, err
	}
	if exactPrice.Float64 == 0 && percentage.Float64 == 0 {
		return nil, shared.ErrPercentAndPriceNull
	}
	if exactPrice.Float64 == 0 {
		exactPrice.Valid = false
	}
	if percentage.Float64 == 0 {
		percentage.Valid = false
	}

	promo := &model.Promotion{
		Name:         payload.Name,
		ExactPrice:   exactPrice,
		Percentage:   percentage,
		MinimumSpend: decimal.NewFromFloat(payload.MinimumSpend),
		Quota:        payload.Quota,
		ShopID:       shop.ID,
		StartedAt:    startedAt,
		ExpiredAt:    expiredAt,
	}
	return promo, nil
}

func (uc *promotionUsecase) DuplicateShopPromotion(ctx context.Context, promoID int64, sellerID int64) error {
	shop, err := uc.sr.FirstShopById(ctx, int(sellerID))
	if err != nil {
		return err
	}
	promo, err := uc.pr.FirstByID(ctx, promoID)
	if err != nil {
		return err
	}
	if promo == nil {
		return shared.ErrPromoNotFound
	}
	if promo.ShopID != shop.ID {
		return shared.ErrUnauthorizedUser
	}
	promo.Name += " - copy"
	return uc.pr.Create(ctx, promo)
}

func (uc *promotionUsecase) DeleteShopPromotion(ctx context.Context, promoID int64, sellerID int64) error {
	shop, err := uc.sr.FirstShopById(ctx, int(sellerID))
	if err != nil {
		return err
	}
	promo, err := uc.pr.FirstByID(ctx, promoID)
	if err != nil {
		return err
	}
	if promo == nil {
		return shared.ErrPromoNotFound
	}
	if promo.ShopID != shop.ID {
		return shared.ErrUnauthorizedUser
	}
	return uc.pr.SoftDelete(ctx, promoID)
}

func NewPromotionRepository(pr repository.PromotionRepository, sr repository.ShopRepository) PromotionUsecase {
	return &promotionUsecase{
		pr: pr,
		sr: sr,
	}
}
