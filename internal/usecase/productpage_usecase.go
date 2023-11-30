package usecase

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/shopspring/decimal"
)

type (
	ProductPageUsecase interface {
		GetProductDetail(ctx context.Context, productCode string, userID int64) (*dto.ProductPageResponse, error)
	}
	productPageUsecase struct {
		pr  repository.ProductRepository
		sr  repository.ShopRepository
		pvr repository.ProductVariantRepository
		pmr repository.ProductMediaRepository
		vtr repository.VariantTypeRepository
		vgr repository.VariantGroupRepository
		wr  repository.WishlistRepository
		rr  repository.ReviewRepository
		odr repository.OrderDetailRepository
	}
)

func (uc *productPageUsecase) GetProductDetail(ctx context.Context, productCode string, userID int64) (*dto.ProductPageResponse, error) {
	prod, err := uc.pr.FirstProductByCode(ctx, productCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrProductNotFound
		}

		return nil, err
	}

	product := &dto.ProductPageProductDetail{
		Name:        prod.Name,
		Description: prod.Description,
		Weight:      prod.Weight,
	}

	shop, err := uc.sr.FirstShopDetailByAccountID(ctx, prod.SellerID)
	if err != nil {
		return nil, err
	}

	pVars, err := uc.pvr.FindProductVariantByProductID(ctx, prod.ID)
	if err != nil {
		return nil, err
	}
	productVariant := make([]dto.ProductPageProductVariant, 0)
	for _, pVar := range pVars {
		disc := decimal.NewFromFloat32((100.0 - pVar.Discount) / 100.0)
		discountedPrice := pVar.Price.Mul(disc)
		v := dto.ProductPageProductVariant{
			ID:              pVar.ID,
			Price:           pVar.Price.InexactFloat64(),
			Stock:           pVar.Stock,
			Discount:        pVar.Discount,
			DiscountedPrice: discountedPrice.InexactFloat64(),
			VariantType1ID:  pVar.VariantType1ID,
			VariantType2ID:  pVar.VariantType2ID,
		}
		productVariant = append(productVariant, v)
	}

	pMeds, err := uc.pmr.FindProductMediaByProductID(ctx, prod.ID)
	if err != nil {
		return nil, err
	}
	productMedias := []dto.ProductPageProductMedia{{
		MediaUrl:  prod.ThumbnailUrl,
		MediaType: "image",
	}}
	for _, pMed := range pMeds {
		m := dto.ProductPageProductMedia{
			MediaUrl:  pMed.MediaUrl,
			MediaType: pMed.MediaType,
		}
		productMedias = append(productMedias, m)
	}

	vGroupTypes, err := uc.vgr.FindVariantGroupWIthVariantTypeByProductID(ctx, prod.ID)
	if err != nil {
		return nil, err
	}

	isVariant := false
	variants := make([]dto.ProductPageVariantsResponse, 0)
	currentGroup := ""
	for _, vGroupType := range vGroupTypes {
		if vGroupType.GroupName != currentGroup || currentGroup == "default" {
			vTypes := make([]dto.ProductPageVariantType, 0)
			vTypes = append(vTypes, dto.ProductPageVariantType{
				TypeID:   vGroupType.TypeID,
				TypeName: vGroupType.TypeName,
			})
			v := dto.ProductPageVariantsResponse{
				GroupName:    vGroupType.GroupName,
				VariantTypes: vTypes,
			}
			variants = append(variants, v)
			currentGroup = vGroupType.GroupName
			continue
		}
		lastIdx := len(variants) - 1
		variants[lastIdx].VariantTypes = append(variants[lastIdx].VariantTypes, dto.ProductPageVariantType{
			TypeID:   vGroupType.TypeID,
			TypeName: vGroupType.TypeName,
		})
	}

	isVariant = variants[0].GroupName == constant.ProductVariantDefault ||
		variants[1].GroupName == constant.ProductVariantDefault

	lowestPrice := productVariant[0].Price
	highestPrice := productVariant[0].Price

	for _, pv := range productVariant {
		if pv.Price < lowestPrice {
			lowestPrice = pv.Price
		}

		if pv.Price > highestPrice {
			highestPrice = pv.Price
		}
	}

	wishlist, err := uc.wr.FirstByUserAndProduct(ctx, userID, prod.ID)
	if err != nil {
		return nil, err
	}

	wishlistCtr, err := uc.wr.CountByProductID(ctx, prod.ID)
	if err != nil {
		return nil, err
	}

	rate, err := uc.rr.RateOfProduct(ctx, productCode)
	if err != nil {
		return nil, err
	}
	rating := rate.RateSum / rate.RateCount
	if math.IsNaN(rating) {
		rating = 0
	}

	ratingCount, err := uc.rr.CountRatingByProductID(ctx, productCode)
	if err != nil {
		return nil, err
	}

	totalSold, err := uc.odr.CountOrderByProductCode(ctx, productCode)
	if err != nil {
		return nil, err
	}

	res := &dto.ProductPageResponse{
		Product:         product,
		Shop:            shop,
		ProductVariants: productVariant,
		ProductMedias:   productMedias,
		VariantGroup1:   variants[0],
		VariantGroup2:   variants[1],
		IsVariant:       isVariant,
		HighPrice:       highestPrice,
		LowPrice:        lowestPrice,
		WishlistCtr:     wishlistCtr.Counter,
		IsInWishlist:    wishlist != nil,
		Rating:          shared.RoundFloat(rating, 1),
		ReviewCount:     *ratingCount,
		TotalSold:       *totalSold,
	}

	return res, nil
}

func NewProductPageUsecase(
	pr repository.ProductRepository,
	sr repository.ShopRepository,
	pvr repository.ProductVariantRepository,
	pmr repository.ProductMediaRepository,
	vtr repository.VariantTypeRepository,
	vgr repository.VariantGroupRepository,
	wr repository.WishlistRepository,
	rr repository.ReviewRepository,
	odr repository.OrderDetailRepository,
) ProductPageUsecase {
	return &productPageUsecase{
		pr:  pr,
		pvr: pvr,
		sr:  sr,
		pmr: pmr,
		vtr: vtr,
		vgr: vgr,
		wr:  wr,
		rr:  rr,
		odr: odr,
	}
}
