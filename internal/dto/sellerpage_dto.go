package dto

type (
	GetSellerDetail struct {
		SellerId     int64  `db:"seller_id"`
		ShopName     string `db:"shop_name"`
		ProductCount string `db:"product_counts"`
		Years        int    `db:"years"`
	}
	GetSellerDetailResponseBody struct {
		ShopName     string                            `json:"shop_name"`
		ProductCount string                            `json:"product_counts"`
		Years        string                            `json:"years"`
		Categories   []string                          `json:"categories"`
		BestSeller   []SearchSellerProductResponseItem `json:"best_seller"`
		Products     []SearchSellerProductResponseItem `json:"products"`
		Pagination   SearchSellerProductResponse       `json:"pagination"`
	}
)

type (
	SearchSellerProductPayload struct {
		SearchTerm   string `validate:"omitempty"`
		Page         int    `validate:"omitempty,numeric,gte=1"`
		SortBy       string `validate:"omitempty,oneof=created_at price most_purchased"`
		SortDesc     bool   `validate:"omitempty,boolean"`
		CategoryName string `validate:"omitempty"`
	}
	CountSellerProductBySearchTermPayload struct {
		SearchTerm string `validate:"omitempty"`
	}
	SearchSellerProductResponseDb struct {
		ProductCode    string  `db:"product_code"`
		ProductName    string  `db:"product_name"`
		ThumbnailURL   string  `db:"thumbnail_url"`
		BasePrice      float64 `db:"base_price"`
		DiscountPrice  float64 `db:"discount_price"`
		Discount       float32 `db:"discount"`
		DistrictName   string  `db:"district_name"`
		CountPurchased int64   `db:"count_purchased"`
	}
	SearchSellerProductResponseItem struct {
		ProductCode    string  `json:"product_code"`
		ProductName    string  `json:"product_name"`
		ThumbnailURL   string  `json:"thumbnail_url"`
		BasePrice      float64 `json:"base_price"`
		DiscountPrice  float64 `json:"discount_price"`
		Discount       float32 `json:"discount"`
		DistrictName   string  `json:"district_name"`
		Rating         float64 `json:"rating"`
		CountPurchased int64   `json:"count_purchased"`
	}
	SearchSellerProductResponse struct {
		Page         int    `json:"page"`
		TotalPage    int    `json:"total_page"`
		TotalProduct int    `json:"total_product"`
		SearchTerm   string `json:"search"`
	}
)
