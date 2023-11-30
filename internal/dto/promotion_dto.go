package dto

import (
	"database/sql"

	"github.com/lil-oren/rest/internal/constant"
)

type (
	PromotionByShopParams struct {
		Status constant.PromotionStatusType `form:"status" validate:"omitempty,oneof=ONGOING COMING ENDED"`
		Page   int                          `form:"page" validate:"required,gt=0"`
	}
	PromotionByShopResponseItems struct {
		ID           int64   `json:"id"`
		Name         string  `json:"name"`
		ExactPrice   float64 `json:"exact_price"`
		Percentage   float64 `json:"Percentage"`
		MinimumSpend float64 `json:"minimum_spend"`
		Quota        int     `json:"quota"`
		StartedAt    string  `json:"started_at"`
		ExpiredAt    string  `json:"expired_at"`
	}
	PromotionByShopResponse struct {
		Items       []PromotionByShopResponseItems `json:"items"`
		TotalData   int                            `json:"total_data"`
		TotalPage   int                            `json:"total_page"`
		CurrentPage int                            `json:"current_page"`
	}
)

type (
	UpsertShopPromotionRequestBody struct {
		Name         string  `json:"name" validate:"required"`
		ExactPrice   float64 `json:"exact_price" validate:"omitempty,gt=0"`
		Percentage   float64 `json:"percentage" validate:"omitempty,gt=0,lte=100"`
		MinimumSpend float64 `json:"minimum_spend" validate:"required,gte=0"`
		Quota        int     `json:"quota" validate:"required,gt=0"`
		StartedAt    string  `json:"started_at" validate:"required"`
		ExpiredAt    string  `json:"expired_at" validate:"required"`
	}
	UpsertShopPromotionPayload struct {
		SellerID     int64
		Name         string
		ExactPrice   float64
		Percentage   float64
		MinimumSpend float64
		Quota        int
		StartedAt    string
		ExpiredAt    string
	}
)

type (
	PromotionDetail struct {
		PromotionName string          `db:"promotion_name"`
		ExactPrice    sql.NullFloat64 `db:"exact_price"`
		Percentage    sql.NullFloat64 `db:"percentage"`
		MinimumSpend  float64         `db:"minimum_spend"`
	}
)
