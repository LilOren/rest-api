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
	OrderUsecase interface {
		GetAllOrder(ctx context.Context, accountId int64, params *dto.OrderParams) ([]dto.OrderList, error)
		GetAllOrderMetadata(ctx context.Context, accountId int64, params *dto.OrderParams) (*dto.OrderPage, error)
		CreateOrder(ctx context.Context, accountId int, payload dto.CreateOrderRequestPayload) error
		CancelOrder(ctx context.Context, orderId int64, userId int64) error
		ReceiveOrder(ctx context.Context, orderId int64, userId int64) error
	}
	orderUsecase struct {
		or  repository.OrderRepository
		cr  repository.CartRepository
		er  repository.WalletRepository
		ror repository.RajaOngkirRepository
		aar repository.AccountAddressRepository
		cor repository.CourierRepository
		pr  repository.ProductRepository
		tr  repository.TransactionRepository
		prr repository.PromotionRepository
	}
)

func (ou *orderUsecase) GetAllOrder(ctx context.Context, accountId int64, params *dto.OrderParams) ([]dto.OrderList, error) {
	orders, err := ou.or.FindOrderByBuyerID(ctx, accountId, params)
	if err != nil {
		return nil, shared.ErrFindOrder
	}

	orderRes := make([]dto.OrderList, 0)
	currentId := int64(0)
	for i, order := range orders {
		if order.ID != currentId {
			products := make([]dto.OrderProducts, 0)
			products = append(products, dto.OrderProducts{
				ProductCode:   order.ProductCode,
				ProductName:   order.ProductName,
				ThumbnailUrl:  order.ThumbnailUrl,
				VariantName:   order.VariantName,
				Quantity:      order.Quantity,
				SubTotalPrice: order.SubTotalPrice.InexactFloat64(),
			})
			v := dto.OrderList{
				ID:                  order.ID,
				Status:              order.Status,
				ShopName:            order.ShopName,
				Products:            products,
				ReceiverName:        order.ReceiverName,
				ReceiverPhoneNumber: order.ReceiverPhoneNumber,
				Address:             order.Address,
				CourierName:         order.CourierName,
				DeliveryCost:        order.DeliveryCost.InexactFloat64(),
				ETA:                 order.ETA.Time.Format("2006-01-02 15:04:05"),
				TotalPrice:          order.SubTotalPrice.InexactFloat64() + order.DeliveryCost.InexactFloat64() - order.PromotionAmount.Float64,
			}
			orderRes = append(orderRes, v)
			currentId = order.ID
			continue
		}
		lastIdx := len(orderRes) - 1
		orderRes[lastIdx].Products = append(orderRes[lastIdx].Products, dto.OrderProducts{
			ProductCode:   order.ProductCode,
			ProductName:   order.ProductName,
			ThumbnailUrl:  order.ThumbnailUrl,
			VariantName:   order.VariantName,
			Quantity:      order.Quantity,
			SubTotalPrice: order.SubTotalPrice.InexactFloat64(),
		})
		orderRes[lastIdx].TotalPrice += order.SubTotalPrice.InexactFloat64()
		if i == len(orders) {
			v := dto.OrderList{
				ID:                  order.ID,
				Status:              order.Status,
				ShopName:            order.ShopName,
				Products:            orderRes[lastIdx].Products,
				ReceiverName:        order.ReceiverName,
				ReceiverPhoneNumber: order.ReceiverPhoneNumber,
				Address:             order.Address,
				CourierName:         order.CourierName,
				DeliveryCost:        order.DeliveryCost.InexactFloat64(),
				ETA:                 order.ETA.Time.Format("2006-01-02 15:04:05"),
				TotalPrice:          order.SubTotalPrice.InexactFloat64() + order.DeliveryCost.InexactFloat64() - order.PromotionAmount.Float64,
			}
			orderRes = append(orderRes, v)
		}
	}
	return orderRes, nil
}

func (ou *orderUsecase) GetAllOrderMetadata(ctx context.Context, accountId int64, params *dto.OrderParams) (*dto.OrderPage, error) {
	orders, err := ou.or.FindOrderByBuyerIDMetadata(ctx, accountId, params)
	if err != nil {
		return nil, err
	}
	totalData := len(orders)
	totalPage := math.Ceil(float64(totalData) / constant.OrderBuyerDefaultItems)
	res := &dto.OrderPage{
		TotalPage: int(totalPage),
	}
	return res, nil
}

func (ou *orderUsecase) CreateOrder(ctx context.Context, accountId int, payload dto.CreateOrderRequestPayload) error {
	totalPrice := decimal.NewFromFloat32(0)
	totalPricePerOrder := decimal.NewFromFloat32(0)
	totalWeightPerOrder := 0
	priceList := make([]decimal.Decimal, 0)
	costList := make([]decimal.Decimal, 0)
	pricePerOrder := make([]decimal.Decimal, 0)
	promotionAmount := make([]float64, 0)
	promotionName := make([]string, 0)
	servicePriceDec := decimal.NewFromFloat(constant.ServicePrice)
	createTxList := make([]*model.Transaction, 0)

	buyerAddress, err := ou.aar.FirstByID(ctx, int64(payload.BuyerAddressId))
	if err != nil {
		return err
	}

	for _, order := range payload.Orders {
		delivCost := float64(0)
		cart, err := ou.cr.FindCheckedCartByShopID(ctx, int64(accountId), int64(order.ShopId))
		if err != nil {
			return shared.ErrFindCart
		}

		if len(cart) == 0 {
			return shared.ErrNoCheckedCart
		}

		shopAddress, err := ou.aar.FirstShopAddressByShopID(ctx, int64(order.ShopId))
		if err != nil {
			return err
		}

		courier, err := ou.cor.FirstByShopCourierID(ctx, int64(order.CourierId))
		if err != nil {
			return err
		}

		for _, checkedCart := range cart {
			disc := decimal.NewFromFloat32(float32((100 - checkedCart.Discount) / 100))
			qty := decimal.NewFromInt(int64(checkedCart.Qty))
			p, err := ou.pr.FirstProductDetail(ctx, checkedCart.ProductID)
			if err != nil {
				return err
			}

			totalWeightPerOrder += p.Weight * checkedCart.Qty

			price := (checkedCart.BasePrice.Mul(disc)).Mul(qty)
			totalPricePerOrder = totalPricePerOrder.Add(price)
			priceList = append(priceList, price)
		}

		promotionDetail, err := ou.prr.FirstPromotionByShopID(ctx, int64(order.ShopId), int64(order.PromotionId))
		if err != nil {
			if err == sql.ErrNoRows {
				return shared.ErrPromotionNotFound
			}
			return shared.ErrFindPromotion
		}

		if promotionDetail.Percentage.Float64 == 0 && promotionDetail.ExactPrice.Float64 == 0 {
			amount := 0.0
			name := ""
			promotionAmount = append(promotionAmount, amount)
			promotionName = append(promotionName, name)
		}

		if promotionDetail.Percentage.Float64 != 0 && totalPricePerOrder.InexactFloat64() >= promotionDetail.MinimumSpend {
			amount := ((promotionDetail.Percentage.Float64) / 100) * totalPricePerOrder.InexactFloat64()
			name := promotionDetail.PromotionName
			pAmountDec := decimal.NewFromFloat(amount)
			totalPricePerOrder = totalPricePerOrder.Sub(pAmountDec)
			promotionAmount = append(promotionAmount, amount)
			promotionName = append(promotionName, name)
		}

		if promotionDetail.ExactPrice.Float64 != 0 && totalPricePerOrder.InexactFloat64() >= promotionDetail.MinimumSpend {
			amount := promotionDetail.ExactPrice.Float64
			name := promotionDetail.PromotionName
			pAmountDec := decimal.NewFromFloat(amount)
			totalPricePerOrder = totalPricePerOrder.Sub(pAmountDec)
			promotionAmount = append(promotionAmount, amount)
			promotionName = append(promotionName, name)
		}

		roQuery := dto.RajaOngkirGetCostHTTPQueries{
			OriginCityID:      int(shopAddress.DistrictId),
			DestinationCityID: int(buyerAddress.DistrictId),
			Weight:            totalWeightPerOrder,
			CourierCode:       courier.Code,
			CourierService:    courier.ServiceName,
		}

		if roQuery.OriginCityID == roQuery.DestinationCityID && roQuery.CourierCode == "jne" {
			roQuery = shared.ServiceChanger(roQuery)
		}

		cost, err := ou.ror.GetCost(ctx, roQuery)
		if err != nil {
			return err
		}

		delivCost = *cost
		costDec := decimal.NewFromFloat(delivCost)
		costList = append(costList, costDec)
		pricePerOrder = append(pricePerOrder, totalPricePerOrder.Add(costDec))
		totalPrice = totalPrice.Add(totalPricePerOrder).Add(costDec)
		totalPricePerOrder = decimal.NewFromFloat32(0)
		totalWeightPerOrder = 0
	}

	walletUser, err := ou.er.FirstActiveWalletByAccountID(ctx, int64(accountId), constant.UserWalletType)
	if err != nil {
		return shared.ErrFindWallet
	}

	walletTemp, err := ou.er.FirstActiveWalletByAccountID(ctx, int64(accountId), constant.TempWalletType)
	if err != nil {
		return shared.ErrFindWallet
	}

	for _, price := range pricePerOrder {
		createTx := &model.Transaction{
			Amount:       price,
			Title:        constant.PaymentOrderTitle,
			ToWalletID:   walletTemp.ID,
			FromWalletID: sql.NullInt64{Int64: walletUser.ID},
		}
		createTxList = append(createTxList, createTx)
	}

	if (totalPrice.Add(servicePriceDec)).GreaterThan(walletUser.Balance) {
		return shared.ErrInsufficientBalance
	}

	err = ou.or.CreateOrder(ctx, int64(accountId), priceList, costList, promotionAmount, promotionName, payload, createTxList)
	if err != nil {
		return err
	}

	return nil
}

func (ou *orderUsecase) CancelOrder(ctx context.Context, orderId int64, userId int64) error {
	order, err := ou.or.FirstOrderByOrderID(ctx, orderId)
	if err != nil {
		return err
	}
	if order.BuyerId != userId {
		return shared.ErrUnauthorizedUser
	}
	return CancelAndRejectOrder(ctx, order, ou.tr, ou.er, ou.or)
}

func (ou *orderUsecase) ReceiveOrder(ctx context.Context, orderId int64, userId int64) error {
	order, err := ou.or.FirstOrderByOrderID(ctx, orderId)
	if err != nil {
		return err
	}
	if order.BuyerId != userId {
		return shared.ErrUnauthorizedUser
	}
	if string(constant.ArriveOrderStatus) != order.Status {
		return shared.ErrWrongInitialStatus
	}
	transaction, err := ou.tr.FirstTransactionByID(ctx, order.TransactionId)
	if err != nil {
		return err
	}
	shopWallet, err := ou.er.FirstActiveWalletByAccountID(ctx, order.SellerId, constant.ShopWalletType)
	if err != nil {
		return err
	}
	receiveTx := &model.Transaction{
		Amount:       transaction.Amount,
		Title:        constant.TransferTitle,
		FromWalletID: sql.NullInt64{Int64: transaction.ToWalletID, Valid: true},
		ToWalletID:   shopWallet.ID,
	}
	err = ou.or.UpdateReceiveOrder(ctx, order.ID, order.BuyerId, order.SellerId, receiveTx)
	if err != nil {
		if errors.Is(err, shared.ErrUpdateInactiveWallet) {
			return shared.ErrWalletNotActivated
		}
		return err
	}
	return nil
}

func CancelAndRejectOrder(ctx context.Context, order *model.Order,
	tr repository.TransactionRepository,
	er repository.WalletRepository,
	or repository.OrderRepository) error {

	if string(constant.NewOrderStatus) != order.Status {
		return shared.ErrWrongInitialStatus
	}
	transaction, err := tr.FirstTransactionByID(ctx, order.TransactionId)
	if err != nil {
		return err
	}
	cancelTx := &model.Transaction{
		Amount:       transaction.Amount,
		Title:        constant.RefundTitle,
		ToWalletID:   transaction.FromWalletID.Int64,
		FromWalletID: sql.NullInt64{Int64: transaction.ToWalletID},
	}
	err = or.UpdateCancelOrder(ctx, order.ID, order.BuyerId, cancelTx)
	if err != nil {
		if errors.Is(err, shared.ErrUpdateInactiveWallet) {
			return shared.ErrWalletNotActivated
		}
		return err
	}
	return nil
}

func NewOrderUsecase(
	or repository.OrderRepository,
	cr repository.CartRepository,
	er repository.WalletRepository,
	ror repository.RajaOngkirRepository,
	aar repository.AccountAddressRepository,
	cor repository.CourierRepository,
	pr repository.ProductRepository,
	tr repository.TransactionRepository,
	prr repository.PromotionRepository,
) OrderUsecase {
	return &orderUsecase{
		or:  or,
		cr:  cr,
		er:  er,
		ror: ror,
		aar: aar,
		cor: cor,
		pr:  pr,
		tr:  tr,
		prr: prr,
	}
}
