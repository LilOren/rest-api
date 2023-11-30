package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/shopspring/decimal"
)

type (
	CheckoutUsecase interface {
		CalculateCheckoutSummary(ctx context.Context, payload dto.CalculateCheckoutSummaryPayload) (*dto.CalculateCheckoutSummaryResponse, error)
		ListCheckoutItem(ctx context.Context, payload dto.ListCheckoutItemPayload) (*dto.ListCheckoutItemResponse, error)
	}
	checkoutUsecase struct {
		cr  repository.CartRepository
		ror repository.RajaOngkirRepository
		aar repository.AccountAddressRepository
		cor repository.CourierRepository
		pr  repository.ProductRepository
		dr  repository.DistrictRepository
		scr repository.ShopCourierRepository
		wr  repository.WalletRepository
		prr repository.PromotionRepository
	}
)

// ListCheckoutItem implements CheckoutUsecase.
func (uc *checkoutUsecase) ListCheckoutItem(ctx context.Context, payload dto.ListCheckoutItemPayload) (*dto.ListCheckoutItemResponse, error) {
	carts, err := uc.cr.FindCheckedCartByAccountID(ctx, payload.UserID)
	if err != nil {
		return nil, err
	}

	res := new(dto.ListCheckoutItemResponse)
	res.Checkouts = make([]dto.ListCheckout, 0)
	if len(carts) == 0 {
		return res, nil
	}

	currentShop := ""
	for _, cart := range carts {
		if cart.ShopName != currentShop {
			discount := decimal.NewFromFloat(cart.Discount / 100.0)
			totalPrice := cart.BasePrice.Mul(decimal.NewFromInt(int64(cart.Qty)))
			discountedPrice := totalPrice.Mul(discount)
			totalPrice = totalPrice.Sub(discountedPrice)
			price, _ := totalPrice.Float64()

			product, err := uc.pr.FirstProductDetail(ctx, cart.ProductID)
			if err != nil {
				return nil, err
			}
			totalWeight := product.Weight * cart.Qty

			items := make([]dto.ListCheckoutItem, 0)
			items = append(items, dto.ListCheckoutItem{
				Name:        cart.ProductName,
				ImageURL:    cart.ImageUrl,
				Quantity:    cart.Qty,
				Price:       price,
				TotalWeight: totalWeight,
			})

			shopAddress, err := uc.aar.FirstShopAddressByShopID(ctx, cart.ShopID)
			if err != nil {
				return nil, err
			}

			district, err := uc.dr.FirstByID(ctx, shopAddress.DistrictId)
			if err != nil {
				return nil, err
			}

			checkouts := dto.ListCheckout{
				ShopID:   cart.ShopID,
				ShopCity: district.Name,
				ShopName: cart.ShopName,
				Items:    items,
			}
			res.Checkouts = append(res.Checkouts, checkouts)
			currentShop = cart.ShopName
			continue
		}
		discount := decimal.NewFromFloat(cart.Discount / 100.0)
		totalPrice := cart.BasePrice.Mul(decimal.NewFromInt(int64(cart.Qty)))
		discountedPrice := totalPrice.Mul(discount)
		totalPrice = totalPrice.Sub(discountedPrice)
		price, _ := totalPrice.Float64()

		product, err := uc.pr.FirstProductDetail(ctx, cart.ProductID)
		if err != nil {
			return nil, err
		}
		totalWeight := product.Weight * cart.Qty

		checkoutLength := len(res.Checkouts)
		res.Checkouts[checkoutLength-1].Items = append(res.Checkouts[checkoutLength-1].Items,
			dto.ListCheckoutItem{
				Name:        cart.ProductName,
				ImageURL:    cart.ImageUrl,
				Quantity:    cart.Qty,
				Price:       price,
				TotalWeight: totalWeight,
			},
		)
	}
	for idx, checkout := range res.Checkouts {
		tempPrice := float64(0)
		for _, item := range checkout.Items {
			tempPrice += item.Price
		}
		res.TotalPrice += tempPrice

		promotions, err := uc.prr.FindOnGoingByShopIDUnpaginated(ctx, checkout.ShopID)
		if err != nil {
			return nil, err
		}

		promotionDropdown := make([]dto.PromotionDropdown, 0)
		for _, promotion := range promotions {
			temp := dto.PromotionDropdown{
				PromotionID:  promotion.ID,
				MinimumSpend: promotion.MinimumSpend.InexactFloat64(),
			}

			minimumSpend := promotion.MinimumSpend.InexactFloat64()
			if tempPrice >= minimumSpend {
				temp.IsApplicable = true
			}

			if promotion.Percentage.Valid {
				temp.Percentage = promotion.Percentage.Float64
			}

			if promotion.ExactPrice.Valid {
				temp.PriceCut = promotion.ExactPrice.Float64
			}

			promotionDropdown = append(promotionDropdown, temp)
		}

		res.Checkouts[idx].PromotionDropdown = promotionDropdown
	}

	for idx, checkout := range res.Checkouts {
		res.Checkouts[idx].CourierDropdown = make([]dto.DropdownValue, 0)
		couriers, err := uc.scr.FindAvailableCourierByShopID(ctx, checkout.ShopID)
		if err != nil {
			return nil, err
		}

		for _, courier := range couriers {
			temp := dto.DropdownValue{
				Label: courier.CourierName,
				Value: courier.ShopCourierID,
			}
			checkout.CourierDropdown = append(checkout.CourierDropdown, temp)
		}

		res.Checkouts[idx].CourierDropdown = checkout.CourierDropdown
	}

	wallet, err := uc.wr.FirstActiveWalletByAccountID(ctx, payload.UserID, constant.UserWalletType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res.IsWalletActivated = false
			return res, nil
		} else {
			return nil, err
		}
	}

	res.IsWalletActivated = true
	res.Balance, _ = wallet.Balance.Float64()

	return res, nil
}

// CalculateCheckoutSummary implements CheckoutUsecase.
func (uc *checkoutUsecase) CalculateCheckoutSummary(ctx context.Context, payload dto.CalculateCheckoutSummaryPayload) (*dto.CalculateCheckoutSummaryResponse, error) {
	res := new(dto.CalculateCheckoutSummaryResponse)

	// get buyer district
	buyerAddress, err := uc.aar.FirstByID(ctx, payload.BuyerAddressID)
	if err != nil {
		return nil, err
	}

	res.Orders = make([]dto.CalculateCheckoutSummaryOrder, 0)
	for _, od := range payload.OrderDeliveries {
		cartModelOrders, err := uc.cr.FindCheckedCartByShopID(ctx, payload.BuyerID, od.ShopID)
		if err != nil {
			return nil, err
		}

		promotionDetail, err := uc.prr.FirstPromotionByShopID(ctx, od.ShopID, od.PromotionID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, shared.ErrPromotionNotFound
			}
			return nil, shared.ErrFindPromotion
		}

		temp := dto.CalculateCheckoutSummaryOrder{
			ShopID: od.ShopID,
		}

		totalWeight := 0
		for _, cmo := range cartModelOrders {
			decimalQty := decimal.NewFromInt32(int32(cmo.Qty))
			basePrice := cmo.BasePrice
			decimalDiscount := decimal.NewFromFloat(cmo.Discount / 100)
			total := basePrice.Mul(decimalQty)
			discountedTotal := total.Mul(decimalDiscount)
			total = total.Sub(discountedTotal)
			subtotalProduct, _ := total.Float64()

			p, err := uc.pr.FirstProductDetail(ctx, cmo.ProductID)
			if err != nil {
				return nil, err
			}

			temp.SubTotalProduct += subtotalProduct
			res.TotalProduct += cmo.Qty
			totalWeight += p.Weight * cmo.Qty
		}

		// get shop shopAddress
		shopAddress, err := uc.aar.FirstShopAddressByShopID(ctx, od.ShopID)
		if err != nil {
			return nil, err
		}

		if od.ShopCourierID != nil && totalWeight != 0 {

			couriers, err := uc.scr.FindShopCourierByShopId(ctx, od.ShopID)
			if err != nil {
				return nil, err
			}

			isCurrentShopCourier := false
			for _, c := range couriers {
				if *od.ShopCourierID == c.ShopCourierID {
					if !c.IsAvailable {
						return nil, shared.ErrCourierNotAvailable
					}

					isCurrentShopCourier = true
					break
				}
			}

			if !isCurrentShopCourier {
				return nil, shared.ErrCourierNotBelongToCurrentShop
			}

			courier, err := uc.cor.FirstByShopCourierID(ctx, *od.ShopCourierID)
			if err != nil {
				return nil, err
			}

			roQuery := dto.RajaOngkirGetCostHTTPQueries{
				OriginCityID:      int(shopAddress.DistrictId),
				DestinationCityID: int(buyerAddress.DistrictId),
				Weight:            totalWeight,
				CourierCode:       courier.Code,
				CourierService:    courier.ServiceName,
			}

			if roQuery.OriginCityID == roQuery.DestinationCityID && roQuery.CourierCode == "jne" {
				roQuery = shared.ServiceChanger(roQuery)
			}

			cost, err := uc.ror.GetCost(ctx, roQuery)
			if err != nil {
				return nil, err
			}

			temp.DeliveryCost = *cost
		}

		if promotionDetail.Percentage.Float64 == 0 && promotionDetail.ExactPrice.Float64 == 0 {
			temp.SubTotalPromotion = temp.SubTotalProduct
		}

		if promotionDetail.Percentage.Float64 != 0 && temp.SubTotalProduct >= promotionDetail.MinimumSpend {
			temp.SubTotalPromotion = ((100 - promotionDetail.Percentage.Float64) / 100) * temp.SubTotalProduct
		}

		if promotionDetail.ExactPrice.Float64 != 0 && temp.SubTotalProduct >= promotionDetail.MinimumSpend {
			temp.SubTotalPromotion = temp.SubTotalProduct - promotionDetail.ExactPrice.Float64
		}

		temp.Subtotal = temp.SubTotalPromotion + temp.DeliveryCost

		res.Orders = append(res.Orders, temp)
	}

	res.ServicePrice = constant.ServicePrice
	res.SummaryPrice += constant.ServicePrice
	for _, o := range res.Orders {
		res.TotalDeliveryCost += o.DeliveryCost
		res.TotalShopPrice += o.SubTotalPromotion
		res.SummaryPrice += o.Subtotal
	}

	return res, nil
}

func NewCheckoutUsecase(
	cr repository.CartRepository,
	aar repository.AccountAddressRepository,
	cor repository.CourierRepository,
	pr repository.ProductRepository,
	ror repository.RajaOngkirRepository,
	dr repository.DistrictRepository,
	scr repository.ShopCourierRepository,
	wr repository.WalletRepository,
	prr repository.PromotionRepository,
) CheckoutUsecase {
	return &checkoutUsecase{
		cr:  cr,
		aar: aar,
		cor: cor,
		pr:  pr,
		ror: ror,
		dr:  dr,
		scr: scr,
		wr:  wr,
		prr: prr,
	}
}
