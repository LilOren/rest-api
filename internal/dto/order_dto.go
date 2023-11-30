package dto

import (
	"database/sql"

	"github.com/lil-oren/rest/internal/constant"
	"github.com/shopspring/decimal"
)

type (
	CreateOrderRequestBody struct {
		Orders         []OrderPayload `json:"order_deliveries"`
		BuyerAddressId int            `json:"buyer_address_id"`
	}
	CreateOrderRequestPayload struct {
		Orders         []OrderPayload
		BuyerAddressId int
	}
	OrderPayload struct {
		ShopId      int `json:"shop_id"`
		CourierId   int `json:"shop_courier_id"`
		PromotionId int `json:"promotion_id"`
	}
)

type (
	OrderBuyerModel struct {
		ID                  int64           `db:"id"`
		Status              string          `db:"status"`
		ProductCode         string          `db:"product_code"`
		ProductName         string          `db:"product_name"`
		ThumbnailUrl        string          `db:"thumbnail_url"`
		VariantName         string          `db:"variant_name"`
		Quantity            int             `db:"quantity"`
		SubTotalPrice       decimal.Decimal `db:"sub_total_price"`
		DeliveryCost        decimal.Decimal `db:"delivery_cost"`
		ReceiverName        string          `db:"receiver_name"`
		ReceiverPhoneNumber string          `db:"receiver_phone_number"`
		Address             string          `db:"address_detail"`
		CourierName         string          `db:"courier_name"`
		ETA                 sql.NullTime    `db:"estimated_time_arrival"`
		ShopName            string          `db:"shop_name"`
		PromotionAmount     sql.NullFloat64 `db:"promotion_amount"`
	}
	OrderBuyerResponse struct {
		Orders     []OrderList `json:"order"`
		Pagination OrderPage   `json:"pagination"`
	}
)

type (
	OrderParams struct {
		Status constant.OrderStatusType `form:"status" validate:"omitempty,oneof=NEW PROCESS DELIVER ARRIVE RECEIVE CANCEL"`
		Page   int                      `form:"page" validate:"required,gt=0"`
	}
	OrderList struct {
		ID                  int64           `json:"id"`
		Status              string          `json:"status"`
		ShopName            string          `json:"shop_name"`
		Products            []OrderProducts `json:"products"`
		ReceiverName        string          `json:"receiver_name"`
		ReceiverPhoneNumber string          `json:"receiver_phone_number"`
		Address             string          `json:"address_detail"`
		CourierName         string          `json:"courier_name"`
		DeliveryCost        float64         `json:"delivery_cost"`
		ETA                 string          `json:"eta"`
		TotalPrice          float64         `json:"total_price"`
	}
	OrderProducts struct {
		ProductCode   string  `json:"product_code"`
		ProductName   string  `json:"product_name"`
		ThumbnailUrl  string  `json:"thumbnail_url"`
		VariantName   string  `json:"variant_name"`
		Quantity      int     `json:"quantity"`
		SubTotalPrice float64 `json:"sub_total_price"`
	}
	OrderPage struct {
		TotalPage int `json:"total_page"`
	}
)

type CreateOrderModel struct {
	SellerID         int64           `db:"seller_id"`
	ShopID           int64           `db:"shop_id"`
	ShopName         string          `db:"shop_name"`
	CartID           int64           `db:"cart_id"`
	ProductCode      string          `db:"product_code"`
	ProductName      string          `db:"product_name"`
	ProductID        int64           `db:"product_id"`
	ImageUrl         string          `db:"image_url"`
	ProductVariantID int64           `db:"product_variant_id"`
	BasePrice        decimal.Decimal `db:"base_price"`
	Discount         float64         `db:"discount"`
	Qty              int             `db:"quantity"`
	RemainingQty     int             `db:"remaining_quantity"`
	Variant1Name     string          `db:"variant1_name"`
	Variant2Name     string          `db:"variant2_name"`
	IsChecked        bool            `db:"is_checked"`
}
