package dto

type (
	SearchProductPayload struct {
		SearchTerm  string `validate:"omitempty"`
		Page        int    `validate:"omitempty,numeric,gte=1"`
		SortBy      string `validate:"omitempty,oneof=created_at price most_purchased"`
		SortDesc    bool   `validate:"omitempty,boolean"`
		DistrictIDs string `validate:"omitempty"`
		CategoryID  int64  `validate:"omitempty,numeric"`
		MinPrice    float64
		MaxPrice    float64
	}
	CountProductBySearchTermPayload struct {
		SearchTerm  string `validate:"omitempty"`
		DistrictIDs string `validate:"omitempty"`
		CategoryID  int64  `validate:"omitempty,numeric"`
		MinPrice    float64
		MaxPrice    float64
	}
	SearchProductResponseItem struct {
		ProductCode   string  `json:"product_code" db:"product_code"`
		ProductName   string  `json:"name" db:"product_name"`
		ThumbnailURL  string  `json:"image_url" db:"thumbnail_url"`
		BasePrice     float64 `json:"price" db:"base_price"`
		DiscountPrice float64 `json:"discounted_price" db:"discount_price"`
		Discount      float32 `json:"discount" db:"discount"`
		ShopName      string  `json:"shop_name" db:"shop_name"`
		DistrictName  string  `json:"shop_location" db:"district_name"`
		TotalSold     int64   `json:"total_sold" db:"total_sold"`
		Rating        float64 `json:"rating"`
	}
	SearchProductResponse struct {
		Products     []SearchProductResponseItem `json:"products"`
		Page         int                         `json:"page"`
		TotalPage    int                         `json:"total_page"`
		TotalProduct int                         `json:"total_product"`
		SearchTerm   string                      `json:"search"`
	}
)
