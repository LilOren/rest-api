package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
)

type (
	ReviewRepository interface {
		FindByProductCode(ctx context.Context, code string, params *dto.GetAllReviewParams) ([]dto.GetAllReviewModel, error)
		FindByProductCodeMetadata(ctx context.Context, code string, params *dto.GetAllReviewParams) ([]model.Review, error)
		Create(ctx context.Context, review *model.Review, mediaUrl []string) error
		RateOfProduct(ctx context.Context, code string) (*dto.RateOfProductModel, error)
		CountRatingByProductID(ctx context.Context, productCode string) (*int, error)
	}
	reviewRepository struct {
		db *sqlx.DB
	}
)

// CountRatingByProductID implements ReviewRepository.
func (r *reviewRepository) CountRatingByProductID(ctx context.Context, productCode string) (*int, error) {
	qs := `
	SELECT
		COUNT(1)
	FROM
		reviews r
	WHERE
		r.product_code = $1
	`

	row := r.db.QueryRowxContext(ctx, qs, productCode)
	if err := row.Err(); err != nil {
		return nil, err
	}

	count := new(int)
	if err := row.Scan(count); err != nil {
		return nil, err
	}

	return count, nil
}

func (r *reviewRepository) FindByProductCode(ctx context.Context, code string, params *dto.GetAllReviewParams) ([]dto.GetAllReviewModel, error) {
	reviews := make([]dto.GetAllReviewModel, 0)
	query := `
		SELECT 
			r.rating, 
			r."comment",
			a.id AS account_id, 
			a.username,
			rm.image_urls,
			r.created_at 
		FROM reviews r 
		LEFT JOIN (
			SELECT rm.review_id, array_agg(rm.image_url) AS image_urls
			FROM review_medias rm
			GROUP BY rm.review_id
		) rm ON rm.review_id = r.id
		LEFT JOIN accounts a ON r.account_id = a.id 
		WHERE r.product_code = $1 `
	if params.Rate != 0 {
		query += fmt.Sprintf("AND r.rating = %d ", params.Rate)
	}
	switch params.Type {
	case "comment":
		query += "AND rm.image_urls IS NULL "
	case "image":
		query += "AND rm.image_urls IS NOT NULL "
	}
	query += "ORDER BY r.created_at "
	if params.Sort == "desc" || params.Sort == "" {
		query += "DESC "
	}
	query += "LIMIT $2 OFFSET $3"

	limit := constant.ProductReviewDefaultItems
	offset := (params.Page - 1) * limit
	rows, err := r.db.QueryContext(ctx, query, code, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		review := new(dto.GetAllReviewModel)
		if err := rows.Scan(&review.Rating, &review.Comment, &review.AccountID, &review.Username, (*pq.StringArray)(&review.ImageUrls), &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, *review)
	}
	return reviews, nil
}

func (r *reviewRepository) FindByProductCodeMetadata(ctx context.Context, code string, params *dto.GetAllReviewParams) ([]model.Review, error) {
	reviews := make([]model.Review, 0)
	query := `
		SELECT 
			r.*
		FROM reviews r 
		LEFT JOIN (
			SELECT rm.review_id, array_agg(rm.image_url) AS image_urls
			FROM review_medias rm
			GROUP BY rm.review_id
		) rm ON rm.review_id = r.id
		LEFT JOIN accounts a ON r.account_id = a.id 
		WHERE r.product_code = $1 `
	if params.Rate != 0 {
		query += fmt.Sprintf("AND r.rating = %d ", params.Rate)
	}
	switch params.Type {
	case "comment":
		query += "AND rm.image_urls IS NULL "
	case "image":
		query += "AND rm.image_urls IS NOT NULL "
	}
	if err := r.db.SelectContext(ctx, &reviews, query, code); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) Create(ctx context.Context, review *model.Review, mediaUrl []string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query1 := `
		INSERT INTO reviews(rating, comment, account_id, product_code)
		VALUES 
		($1, $2, $3, $4)
		RETURNING id
		`
	reviewID := new(int64)
	if err := tx.Get(reviewID, query1, review.Rating, review.Comment, review.AccountID, review.ProductCode); err != nil {
		return err
	}

	if len(mediaUrl) != 0 {
		rMedias := make([]model.ReviewMedia, 0)
		for _, url := range mediaUrl {
			rMedias = append(rMedias, model.ReviewMedia{
				ReviewID: *reviewID,
				ImageUrl: url,
			})
		}

		query2 := `
			INSERT INTO review_medias(review_id, image_url)
			VALUES 
			(:review_id, :image_url)
		`
		if _, err := tx.NamedExec(query2, rMedias); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *reviewRepository) RateOfProduct(ctx context.Context, code string) (*dto.RateOfProductModel, error) {
	rate := new(dto.RateOfProductModel)
	query := `
		SELECT 
			count(id) AS rate_count,
			COALESCE(sum(rating), 0) AS rate_sum
		FROM reviews r
		WHERE r.product_code = $1
	`
	if err := r.db.GetContext(ctx, rate, query, code); err != nil {
		return nil, err
	}
	return rate, nil
}

func NewReviewRepository(db *sqlx.DB) ReviewRepository {
	return &reviewRepository{
		db: db,
	}
}
