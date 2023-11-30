package usecase

import (
	"context"
	"math"
	"time"

	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
)

type (
	OrderSellerUsecase interface {
		GetAllOrderOfSeller(ctx context.Context, sellerId int64, params *dto.OrderSellerParams) ([]dto.OrderSellerData, error)
		GetAllOrderOfSellerMetadata(ctx context.Context, sellerId int64, params *dto.OrderSellerParams) (*dto.OrderSellerMetadata, error)
		UpdateOrderStatus(ctx context.Context, orderId int64, req *dto.OrderSellerStatusRequest, userId int64) error
		RejectOrder(ctx context.Context, orderId int64, userId int64) error
	}
	orderSellerUsecase struct {
		or repository.OrderRepository
		tr repository.TransactionRepository
		wr repository.WalletRepository
	}
)

func (ouc *orderSellerUsecase) GetAllOrderOfSeller(ctx context.Context, sellerId int64, params *dto.OrderSellerParams) ([]dto.OrderSellerData, error) {
	orders, err := ouc.or.FindOrderBySellerID(ctx, sellerId, params)
	if err != nil {
		return nil, shared.ErrFindOrder
	}
	orderRes := make([]dto.OrderSellerData, 0)
	currentId := int64(0)
	for i, order := range orders {
		if order.ID != currentId {
			products := make([]dto.OrderSellerProducts, 0)
			products = append(products, dto.OrderSellerProducts{
				ProductName:   order.ProductName,
				ThumbnailUrl:  order.ThumbnailUrl,
				VariantName:   order.VariantName,
				SubTotalPrice: order.SubTotalPrice.InexactFloat64(),
				Quantity:      order.Quantity,
			})
			v := dto.OrderSellerData{
				ID:                   order.ID,
				Status:               order.Status,
				Products:             products,
				ReceiverName:         order.ReceiverName,
				ReceiverPhoneNumber:  order.ReceiverPhoneNumber,
				Address:              order.Address,
				CourierName:          order.CourierName,
				ETA:                  order.ETA.Time.Format("2006-01-02"),
				TotalBeforePromotion: order.SubTotalPrice.InexactFloat64(),
				PromotionAmount:      order.PromotionAmount.Float64,
				TotalPrice:           order.SubTotalPrice.InexactFloat64() - order.PromotionAmount.Float64,
			}
			orderRes = append(orderRes, v)
			currentId = order.ID
			continue
		}
		lastIdx := len(orderRes) - 1
		orderRes[lastIdx].Products = append(orderRes[lastIdx].Products, dto.OrderSellerProducts{
			ProductName:   order.ProductName,
			ThumbnailUrl:  order.ThumbnailUrl,
			VariantName:   order.VariantName,
			SubTotalPrice: order.SubTotalPrice.InexactFloat64(),
			Quantity:      order.Quantity,
		})
		orderRes[lastIdx].TotalBeforePromotion += order.SubTotalPrice.InexactFloat64()
		orderRes[lastIdx].TotalPrice += order.SubTotalPrice.InexactFloat64()
		if i == len(orders) {
			v := dto.OrderSellerData{
				ID:                   order.ID,
				Status:               order.Status,
				Products:             orderRes[lastIdx].Products,
				ReceiverName:         order.ReceiverName,
				ReceiverPhoneNumber:  order.ReceiverPhoneNumber,
				Address:              order.Address,
				CourierName:          order.CourierName,
				ETA:                  order.ETA.Time.Format("2006-01-02 15:04:05"),
				TotalBeforePromotion: orderRes[lastIdx].TotalPrice,
				PromotionAmount:      order.PromotionAmount.Float64,
				TotalPrice:           orderRes[lastIdx].TotalPrice - order.PromotionAmount.Float64,
			}
			orderRes = append(orderRes, v)
		}
	}
	return orderRes, nil
}

func (ouc *orderSellerUsecase) GetAllOrderOfSellerMetadata(ctx context.Context, sellerId int64, params *dto.OrderSellerParams) (*dto.OrderSellerMetadata, error) {
	orders, err := ouc.or.FindOrderBySellerIDMetadata(ctx, sellerId, params)
	if err != nil {
		return nil, err
	}
	totalData := len(orders)
	totalPage := math.Ceil(float64(totalData) / constant.OrderSellerDefaultItems)
	res := &dto.OrderSellerMetadata{
		TotalData: totalData,
		TotalPage: int(totalPage),
	}
	return res, nil
}

func (ouc *orderSellerUsecase) UpdateOrderStatus(ctx context.Context, orderId int64, req *dto.OrderSellerStatusRequest, userId int64) error {
	order, err := ouc.or.FirstOrderByOrderID(ctx, orderId)
	if err != nil {
		return err
	}
	if order.SellerId != userId {
		return shared.ErrUnauthorizedUser
	}
	statusMap := map[constant.OrderStatusType]string{
		constant.ProcessOrderStatus: "NEW",
		constant.DeliverOrderStatus: "PROCESS",
		constant.ArriveOrderStatus:  "DELIVER",
	}
	if statusMap[req.NewStatus] != order.Status {
		return shared.ErrWrongInitialStatus
	}
	if req.EstDays < 1 {
		err = ouc.or.UpdateOrderStatus(ctx, orderId, req.NewStatus, nil)
		if err != nil {
			return err
		}
		return nil
	}
	eat := time.Now().Add(time.Hour * (time.Duration(req.EstDays * 24)))
	err = ouc.or.UpdateOrderStatus(ctx, orderId, req.NewStatus, &eat)
	if err != nil {
		return err
	}
	return nil
}

func (ou *orderSellerUsecase) RejectOrder(ctx context.Context, orderId int64, userId int64) error {
	order, err := ou.or.FirstOrderByOrderID(ctx, orderId)
	if err != nil {
		return err
	}
	if order.SellerId != userId {
		return shared.ErrUnauthorizedUser
	}
	return CancelAndRejectOrder(ctx, order, ou.tr, ou.wr, ou.or)
}

func NewOrderSellerUsecase(or repository.OrderRepository, tr repository.TransactionRepository, wr repository.WalletRepository) OrderSellerUsecase {
	return &orderSellerUsecase{
		or: or,
		tr: tr,
		wr: wr,
	}
}
