package dto

import "time"

type (
	AddReviewRequestBody struct {
		ProductCode string   `json:"product_code" validate:"required"`
		Rating      int      `json:"rating" validate:"required,gte=1,lte=5"`
		Comment     string   `json:"comment"`
		ImageUrls   []string `json:"image_urls"`
	}
	AddReviewPayload struct {
		AccountID   int64
		ProductCode string
		Rating      int
		Comment     string
		ImageUrls   []string
	}
)

type (
	GetAllReviewModel struct {
		Rating    int       `db:"rating"`
		Comment   string    `db:"comment"`
		AccountID int64     `db:"account_id"`
		Username  string    `db:"username"`
		ImageUrls []string  `db:"image_urls"`
		CreatedAt time.Time `db:"created_at"`
	}
	GetAllReviewUserResponse struct {
		Rating    int      `json:"rating"`
		Comment   string   `json:"comment"`
		AccountID int64    `json:"account_id"`
		Username  string   `json:"username"`
		ImageUrls []string `json:"image_urls,omitempty"`
		CreatedAt string   `json:"created_at"`
	}
	GetAllReviewResponse struct {
		UserReview  []GetAllReviewUserResponse `json:"user_reviews"`
		TotalReview int                        `json:"total_review"`
		TotalPage   int                        `json:"total_page"`
		CurrentPage int                        `json:"current_page"`
	}
	GetAllReviewParams struct {
		Page int    `form:"page" validate:"gt=0"`
		Rate int    `form:"rate" validate:"omitempty,oneof=1 2 3 4 5"`
		Type string `form:"type" validate:"omitempty,oneof=comment image"`
		Sort string `form:"sort" validate:"omitempty,oneof=asc desc"`
	}
)

type (
	RateOfProductModel struct {
		RateCount float64 `db:"rate_count"`
		RateSum   float64 `db:"rate_sum"`
	}
)
