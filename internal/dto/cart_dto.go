package dto

import "github.com/shopspring/decimal"

type (
	AddToCartRequestBody struct {
		ProductVariantID int64 `json:"product_variant_id" validate:"required"`
		SellerID         int64 `json:"seller_id" validate:"required"`
		Quantity         int   `json:"quantity" validate:"required,gte=1"`
	}
	AddToCartRequestPayload struct {
		ProductVariantID int64
		SellerID         int64
		Quantity         int
	}
)

type (
	IsCheckedCartRequestBody struct {
		IsCheckCarts []IsCheckedCartItem `json:"is_checked_carts" validate:"required"`
	}
	IsCheckedCartResponse struct {
		TotalBasePrice     float64 `json:"total_base_price"`
		TotalDiscountPrice float64 `json:"total_discount_price"`
		TotalPrice         float64 `json:"total_price"`
	}
	IsCheckedModel struct {
		Quantity int             `db:"quantity"`
		Price    decimal.Decimal `db:"price"`
		Discount float64         `db:"discount"`
	}
)

type CartPageModel struct {
	ShopName     string          `db:"shop_name"`
	ShopID       int64           `db:"shop_id"`
	CartID       int64           `db:"cart_id"`
	ProductName  string          `db:"product_name"`
	ProductID    int64           `db:"product_id"`
	ImageUrl     string          `db:"image_url"`
	BasePrice    decimal.Decimal `db:"base_price"`
	Discount     float64         `db:"discount"`
	Qty          int             `db:"quantity"`
	RemainingQty int             `db:"remaining_quantity"`
	Variant1Name string          `db:"variant1_name"`
	Variant2Name string          `db:"variant2_name"`
	IsChecked    bool            `db:"is_checked"`
}

type CartPageResponse struct {
	Items  []CartPageItems       `json:"items"`
	Prices IsCheckedCartResponse `json:"prices"`
}

type CartOrderModel struct {
	SellerID         int64           `db:"seller_id"`
	ShopID           int64           `db:"shop_id"`
	ShopName         string          `db:"shop_name"`
	CartID           int64           `db:"cart_id"`
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

type (
	CartHomeModel struct {
		ProductName  string          `db:"product_name"`
		ThumbnailUrl string          `db:"thumbnail_url"`
		BasePrice    decimal.Decimal `db:"base_price"`
		Discount     float64         `db:"discount"`
		Quantity     int             `db:"quantity"`
	}
	CartHomeResponse struct {
		ProductName  string  `json:"product_name"`
		ThumbnailUrl string  `json:"thumbnail_url"`
		Price        float64 `json:"price"`
		Quantity     int     `json:"quantity"`
	}
)

type (
	CartPageProduct struct {
		CartID        int64   `json:"cart_id"`
		ProductName   string  `json:"product_name"`
		ProductID     int64   `json:"product_id"`
		ImageUrl      string  `json:"image_url"`
		BasePrice     float64 `json:"base_price"`
		DiscountPrice float64 `json:"discount_price"`
		Discount      float64 `json:"discount"`
		Qty           int     `json:"quantity"`
		RemainingQty  int     `json:"remaining_quantity"`
		Variant1Name  string  `json:"variant1_name"`
		Variant2Name  string  `json:"variant2_name"`
		IsChecked     bool    `json:"is_checked"`
	}
	ProductVariantForCartModel struct {
		ID       int64 `db:"id"`
		Stock    int   `db:"stock"`
		SellerID int64 `db:"seller_id"`
	}
	UpdateQuantityRequestBody struct {
		Quantity int64 `json:"quantity" validate:"required,gte=1"`
	}
	IsCheckedCartItem struct {
		CartID    int64 `json:"cart_id" validate:"required"`
		IsChecked bool  `json:"is_checked" validate:"required"`
	}
	CartPageItems struct {
		ShopName string            `json:"seller_name"`
		ShopID   int64             `json:"seller_id"`
		Products []CartPageProduct `json:"products"`
	}
)
