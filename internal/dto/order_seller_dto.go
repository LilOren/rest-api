package dto

import (
	"database/sql"

	"github.com/lil-oren/rest/internal/constant"
	"github.com/shopspring/decimal"
)

type (
	OrderSellerStatusRequest struct {
		NewStatus constant.OrderStatusType
		EstDays   int `json:"est_days" validate:"required,gte=1,lte=3"`
	}
	OrderSellerModel struct {
		ID                  int64           `db:"id"`
		Status              string          `db:"status"`
		ProductCode         string          `db:"product_code"`
		ProductName         string          `db:"product_name"`
		ThumbnailUrl        string          `db:"thumbnail_url"`
		VariantName         string          `db:"variant_name"`
		SubTotalPrice       decimal.Decimal `db:"sub_total_price"`
		Quantity            int             `db:"quantity"`
		ReceiverName        string          `db:"receiver_name"`
		ReceiverPhoneNumber string          `db:"receiver_phone_number"`
		Address             string          `db:"address_detail"`
		CourierName         string          `db:"courier_name"`
		ETA                 sql.NullTime    `db:"estimated_time_arrival"`
		PromotionAmount     sql.NullFloat64 `db:"promotion_amount"`
	}
	OrderSellerResponse struct {
		OrdersData []OrderSellerData `json:"order_data"`
		TotalData  int               `json:"total_data"`
		TotalPage  int               `json:"total_page"`
	}
	OrderSellerData struct {
		ID                   int64                 `json:"id"`
		Status               string                `json:"status"`
		Products             []OrderSellerProducts `json:"products"`
		ReceiverName         string                `json:"receiver_name"`
		ReceiverPhoneNumber  string                `json:"receiver_phone_number"`
		Address              string                `json:"address_detail"`
		CourierName          string                `json:"courier_name"`
		ETA                  string                `json:"eta"`
		TotalBeforePromotion float64               `json:"total_before_promotion"`
		PromotionAmount      float64               `json:"promotion_amount"`
		TotalPrice           float64               `json:"total_price"`
	}
)

type (
	OrderSellerParams struct {
		Status constant.OrderStatusType `form:"status" validate:"omitempty,oneof=NEW PROCESS DELIVER ARRIVE RECEIVE CANCEL"`
		Page   int                      `form:"page" validate:"required,gt=0"`
	}
	OrderSellerProducts struct {
		ProductName   string  `json:"product_name"`
		ThumbnailUrl  string  `json:"thumbnail_url"`
		VariantName   string  `json:"variant_name"`
		SubTotalPrice float64 `json:"sub_total_price"`
		Quantity      int     `json:"quantity"`
	}
	OrderSellerMetadata struct {
		TotalData int
		TotalPage int
	}
)
