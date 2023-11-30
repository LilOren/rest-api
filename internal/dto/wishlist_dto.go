package dto

import "github.com/shopspring/decimal"

type (
	WishlistRequestBody struct {
		ProductCode string `json:"product_code" validate:"required"`
	}
	WishlistPayload struct {
		UserID      int64
		ProductCode string
	}
)

type (
	WishlistParams struct {
		Page int `form:"page" validate:"required,gt=0"`
	}
	WishlistUserModel struct {
		ID           int64           `db:"id"`
		ProductCode  string          `db:"product_code"`
		ProductName  string          `db:"product_name"`
		ThumbnailURL string          `db:"thumbnail_url"`
		BasePrice    decimal.Decimal `db:"base_price"`
		Discount     float32         `db:"discount"`
		ShopName     string          `db:"shop_name"`
		DistrictName string          `db:"district_name"`
	}
	WishlistUserResponse struct {
		Items       []WishlistUserResponseItems `json:"items"`
		CurrentPage int                         `json:"current_page"`
		TotalPage   int                         `json:"total_page"`
		TotalData   int                         `json:"total_data"`
	}
	WishlistUserResponseItems struct {
		ID            int64   `json:"id"`
		ProductCode   string  `json:"product_code"`
		ProductName   string  `json:"product_name"`
		ThumbnailURL  string  `json:"thumbnail_url"`
		BasePrice     float64 `json:"base_price"`
		Discount      float32 `json:"discount"`
		DiscountPrice float64 `json:"discount_price"`
		ShopName      string  `json:"shop_name"`
		DistrictName  string  `json:"district_name"`
		Rating        float64 `json:"rating"`
	}
)

type (
	WishlistCountPayload struct {
		Counter int `db:"counter"`
	}
)
