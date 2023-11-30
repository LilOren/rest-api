package dto

import "github.com/shopspring/decimal"

type (
	HomePageProductModel struct {
		ImageUrl        string          `db:"media_url"`
		ProductCode     string          `db:"product_code"`
		Name            string          `db:"name"`
		Price           decimal.Decimal `db:"price"`
		DiscountedPrice decimal.Decimal `db:"discounted_price"`
		Discount        float32         `db:"discount"`
		TotalSold       int             `db:"total_sold"`
		ShopName        string          `db:"shop_name"`
		ShopLocation    string          `db:"shop_location"`
	}
	HomePageProductResponseBody struct {
		ImageUrl        string  `json:"image_url"`
		ProductCode     string  `json:"product_code"`
		Name            string  `json:"name"`
		Price           float64 `json:"price"`
		DiscountedPrice float64 `json:"discounted_price"`
		Discount        float32 `json:"discount"`
		TotalSold       int     `json:"total_sold"`
		ShopName        string  `json:"shop_name"`
		ShopLocation    string  `json:"shop_location"`
		Rating          float64 `json:"rating"`
	}
	HomePageCategoryPair struct {
		FirstLevelID          int64  `db:"first_level_id"`
		FirstLevelName        string `db:"first_level_name"`
		FirstCategoryImageURL string `db:"image_url"`
		SecondLevelID         int64  `db:"second_level_id"`
	}
	HomePageCategoryResponseBody struct {
		TopLevelCategoryID   int64  `json:"top_category_id"`
		ChildLevelCategoryID int64  `json:"child_category_id"`
		CategoryName         string `json:"category_name"`
		ImageURL             string `json:"image_url"`
	}
)
