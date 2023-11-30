package usecase

import (
	"context"

	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/shopspring/decimal"
)

type (
	CartUsecase interface {
		GetCartPageAllProducts(ctx context.Context, accountId int64) (*dto.CartPageResponse, error)
		AddToCart(ctx context.Context, product *dto.AddToCartRequestPayload, accountId int64) error
		UpdateQuantityItem(ctx context.Context, cartId int64, quantity int) error
		DeleteItem(ctx context.Context, cartId int64) error
		UpdateIsCheckCart(ctx context.Context, items []dto.IsCheckedCartItem) error
		GetTotalPriceChecked(ctx context.Context, accountId int64) (*dto.IsCheckedCartResponse, error)
	}
	cartUsecase struct {
		cr  repository.CartRepository
		pvr repository.ProductVariantRepository
	}
)

func (cuc *cartUsecase) GetCartPageAllProducts(ctx context.Context, accountId int64) (*dto.CartPageResponse, error) {
	cartProducts, err := cuc.cr.FindByAccountID(ctx, accountId)
	if err != nil {
		return nil, err
	}
	cartShop := make([]dto.CartPageItems, 0)
	currentShop := ""
	for _, cartProduct := range cartProducts {
		if cartProduct.ShopName != currentShop {
			pVariants := make([]dto.CartPageProduct, 0)
			pVariants = append(pVariants, dto.CartPageProduct{
				CartID:        cartProduct.CartID,
				ProductName:   cartProduct.ProductName,
				ProductID:     cartProduct.ProductID,
				ImageUrl:      cartProduct.ImageUrl,
				BasePrice:     cartProduct.BasePrice.InexactFloat64(),
				DiscountPrice: cartProduct.BasePrice.Mul(decimal.NewFromFloat(((100 - cartProduct.Discount) / 100))).InexactFloat64(),
				Discount:      cartProduct.Discount,
				Qty:           cartProduct.Qty,
				RemainingQty:  cartProduct.RemainingQty,
				Variant1Name:  cartProduct.Variant1Name,
				Variant2Name:  cartProduct.Variant2Name,
				IsChecked:     cartProduct.IsChecked,
			})
			v := dto.CartPageItems{
				ShopName: cartProduct.ShopName,
				ShopID:   cartProduct.ShopID,
				Products: pVariants,
			}
			cartShop = append(cartShop, v)
			currentShop = cartProduct.ShopName
			continue
		}
		cartLenght := len(cartShop)
		cartShop[cartLenght-1].Products = append(cartShop[cartLenght-1].Products, dto.CartPageProduct{
			CartID:        cartProduct.CartID,
			ProductName:   cartProduct.ProductName,
			ProductID:     cartProduct.ProductID,
			ImageUrl:      cartProduct.ImageUrl,
			BasePrice:     cartProduct.BasePrice.InexactFloat64(),
			DiscountPrice: cartProduct.BasePrice.Mul(decimal.NewFromFloat(((100 - cartProduct.Discount) / 100))).InexactFloat64(),
			Discount:      cartProduct.Discount,
			Qty:           cartProduct.Qty,
			RemainingQty:  cartProduct.RemainingQty,
			Variant1Name:  cartProduct.Variant1Name,
			Variant2Name:  cartProduct.Variant2Name,
			IsChecked:     cartProduct.IsChecked,
		})
	}
	prices, err := cuc.GetTotalPriceChecked(ctx, accountId)
	if err != nil {
		return nil, err
	}
	resp := &dto.CartPageResponse{Items: cartShop, Prices: *prices}
	return resp, nil
}

func (cuc *cartUsecase) AddToCart(ctx context.Context, product *dto.AddToCartRequestPayload, accountId int64) error {
	if product.SellerID == accountId {
		return shared.ErrOwnSellerProduct
	}
	productVariant, err := cuc.pvr.FirstProductVariantByIDForCart(ctx, product.ProductVariantID)
	if err != nil {
		return err
	}
	if product.SellerID != productVariant.SellerID {
		return shared.ErrDifferentSeller
	}
	if product.Quantity > productVariant.Stock {
		return shared.ErrQuantityMoreThanStock
	}
	cartProduct, err := cuc.cr.FirstByProductVariantID(ctx, product.ProductVariantID, accountId)
	if err == nil {
		if cartProduct.Quantity+product.Quantity > productVariant.Stock {
			return shared.ErrQuantityMoreThanStock
		}
		err = cuc.cr.IncreaseQuantityByID(ctx, cartProduct.ID, product.Quantity)
		if err != nil {
			return err
		}
		return nil
	}
	item := &model.Cart{
		ProductVariantId: product.ProductVariantID,
		AccountId:        accountId,
		SellerId:         product.SellerID,
		Quantity:         product.Quantity,
	}
	err = cuc.cr.Create(ctx, item)
	if err != nil {
		return err
	}
	return nil
}

func (cuc *cartUsecase) UpdateQuantityItem(ctx context.Context, cartId int64, quantity int) error {
	product, err := cuc.cr.FirstCart(ctx, cartId)
	if product == nil {
		return shared.ErrCartNotFound
	}
	if err != nil {
		return err
	}
	productVariant, err := cuc.pvr.FirstProductVariantByID(ctx, product.ProductVariantId)
	if err != nil {
		return err
	}
	if quantity > int(productVariant.Stock) {
		return shared.ErrQuantityMoreThanStock
	}
	if err = cuc.cr.UpdateQuantity(ctx, quantity, cartId); err != nil {
		return err
	}
	return nil
}

func (cuc *cartUsecase) DeleteItem(ctx context.Context, cartId int64) error {
	item, err := cuc.cr.FirstCart(ctx, cartId)
	if item == nil {
		return shared.ErrCartNotFound
	}
	if err != nil {
		return err
	}
	if err = cuc.cr.DeleteCart(ctx, cartId); err != nil {
		return err
	}
	return nil
}

func (cuc *cartUsecase) UpdateIsCheckCart(ctx context.Context, items []dto.IsCheckedCartItem) error {
	itemsModel := make([]model.Cart, 0)
	for _, val := range items {
		item := model.Cart{
			ID:        val.CartID,
			IsChecked: val.IsChecked,
		}
		itemsModel = append(itemsModel, item)
	}
	if err := cuc.cr.UpdateCheck(ctx, itemsModel); err != nil {
		return err
	}
	return nil
}

func (cuc *cartUsecase) GetTotalPriceChecked(ctx context.Context, accountId int64) (*dto.IsCheckedCartResponse, error) {
	carts, err := cuc.cr.FindCheckedForPrice(ctx, accountId)
	if err != nil {
		return nil, err
	}
	var basePrice, discountPrice float64 = 0, 0
	for _, val := range carts {
		basePrice += (val.Price.Mul(decimal.NewFromInt(int64(val.Quantity))).InexactFloat64())
		discountPrice += (val.Price.Mul(decimal.NewFromFloat((val.Discount / 100))).Mul(decimal.NewFromFloat(float64(val.Quantity))).InexactFloat64())
	}
	totalPrice := basePrice - discountPrice
	res := &dto.IsCheckedCartResponse{
		TotalBasePrice:     basePrice,
		TotalDiscountPrice: discountPrice,
		TotalPrice:         totalPrice,
	}
	return res, nil
}

func NewCartUsecase(cr repository.CartRepository, pvr repository.ProductVariantRepository) CartUsecase {
	return &cartUsecase{
		cr:  cr,
		pvr: pvr,
	}
}
