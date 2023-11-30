package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
)

type (
	PromotionRepository interface {
		FindByShopID(ctx context.Context, shopID int64, params *dto.PromotionByShopParams) ([]model.Promotion, error)
		FindByShopIDMetadata(ctx context.Context, shopID int64, params *dto.PromotionByShopParams) ([]model.Promotion, error)
		FirstByID(ctx context.Context, promotionID int64) (*model.Promotion, error)
		Create(ctx context.Context, promotion *model.Promotion) error
		FindOnGoingByShopIDUnpaginated(ctx context.Context, shopID int64) ([]model.Promotion, error)
		Update(ctx context.Context, promotion *model.Promotion) error
		FirstPromotionByShopID(ctx context.Context, shopID, promotionID int64) (*dto.PromotionDetail, error)
		SoftDelete(ctx context.Context, promotionID int64) error
	}
	promotionRepository struct {
		db *sqlx.DB
	}
)

// FindOnGoingByShopIDUnpaginated implements PromotionRepository.
func (r *promotionRepository) FindOnGoingByShopIDUnpaginated(ctx context.Context, shopID int64) ([]model.Promotion, error) {
	qs := `
	SELECT 
		* 
	FROM 
		promotions p
	WHERE
		p.shop_id = $1
		AND
		p.started_at <= NOW()
		AND
		p.expired_at > NOW()
		AND
		p.deleted_at IS NULL
	`

	m := make([]model.Promotion, 0)
	if err := r.db.SelectContext(ctx, &m, qs, shopID); err != nil {
		return nil, err
	}

	return m, nil
}

func (r *promotionRepository) FindByShopID(ctx context.Context, shopID int64, params *dto.PromotionByShopParams) ([]model.Promotion, error) {
	promotions := make([]model.Promotion, 0)
	query := `
		SELECT * FROM promotions p WHERE p.shop_id = $1
	`
	if params.Status == constant.OngoingPromotionStatus {
		query += "AND p.expired_at > now() AND p.started_at < now() "
	}
	if params.Status == constant.ComingPromotionStatus {
		query += "AND p.started_at > now() "
	}
	if params.Status == constant.EndedPromotionStatus {
		query += "AND p.expired_at < now() "
	}
	query += "AND p.deleted_at IS NULL ORDER BY p.updated_at DESC LIMIT $2 OFFSET $3"
	limit := constant.PromotionShopDefaultItems
	offset := (params.Page - 1) * limit
	if err := r.db.SelectContext(ctx, &promotions, query, shopID, limit, offset); err != nil {
		return nil, err
	}
	return promotions, nil
}

func (r *promotionRepository) FindByShopIDMetadata(ctx context.Context, shopID int64, params *dto.PromotionByShopParams) ([]model.Promotion, error) {
	promotions := make([]model.Promotion, 0)
	query := `
		SELECT * FROM promotions p WHERE p.shop_id = $1
	`
	if params.Status == constant.OngoingPromotionStatus {
		query += "AND p.expired_at > now() AND p.started_at < now() "
	}
	if params.Status == constant.ComingPromotionStatus {
		query += "AND p.started_at > now() "
	}
	if params.Status == constant.EndedPromotionStatus {
		query += "AND p.expired_at < now() "
	}
	query += "AND p.deleted_at IS NULL"
	if err := r.db.SelectContext(ctx, &promotions, query, shopID); err != nil {
		return nil, err
	}
	return promotions, nil
}

func (r *promotionRepository) FirstByID(ctx context.Context, promotionID int64) (*model.Promotion, error) {
	promotion := new(model.Promotion)
	query := `
		SELECT * FROM promotions p WHERE p.id = $1 AND p.deleted_at IS NULL
	`
	if err := r.db.GetContext(ctx, promotion, query, promotionID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return promotion, nil
}

func (r *promotionRepository) Create(ctx context.Context, promotion *model.Promotion) error {
	query := `
	INSERT INTO promotions (
		name, exact_price, 
		percentage, minimum_spend, 
		quota, shop_id, 
		started_at, expired_at
	) VALUES
	(
		:name, :exact_price, 
		:percentage, :minimum_spend,
		:quota, :shop_id, 
		:started_at, :expired_at
	)
	`
	_, err := r.db.NamedExecContext(ctx, query, promotion)
	if err != nil {
		return err
	}

	return nil
}

func (r *promotionRepository) Update(ctx context.Context, promotion *model.Promotion) error {
	query := `
	UPDATE promotions
	SET
		name = :name,
		exact_price = :exact_price,
		percentage = :percentage,
		minimum_spend = :minimum_spend,
		quota = :quota,
		started_at = :started_at,
		expired_at = :expired_at,
		updated_at = now() 
	WHERE id = :id AND deleted_at IS NULL;
	`
	_, err := r.db.NamedExecContext(ctx, query, promotion)
	if err != nil {
		return err
	}

	return nil
}

func (r *promotionRepository) FirstPromotionByShopID(ctx context.Context, shopID, promotionID int64) (*dto.PromotionDetail, error) {
	if promotionID == 0 {
		return &dto.PromotionDetail{
			PromotionName: "",
			ExactPrice: sql.NullFloat64{
				Float64: 0,
				Valid:   true,
			},
			Percentage: sql.NullFloat64{
				Float64: 0,
				Valid:   true,
			},
			MinimumSpend: 0,
		}, nil
	}
	promotion := new(dto.PromotionDetail)
	qs := `
	SELECT
		ps.name AS promotion_name,
		ps.exact_price,
		ps.percentage,
		ps.minimum_spend
	FROM
		promotions ps
	WHERE
		ps.id = $1 AND
		ps.shop_id = $2 AND 
		ps.quota > 0 AND 
		ps.started_at < now() AND 
		ps.expired_at >= now()
	`

	err := r.db.GetContext(ctx, promotion, qs, promotionID, shopID)
	if err != nil {
		return nil, err
	}

	return promotion, nil
}

func (r *promotionRepository) SoftDelete(ctx context.Context, promotionID int64) error {
	query := `
	UPDATE promotions
	SET
		deleted_at = now() 
	WHERE id = $1 AND deleted_at IS NULL;
	`
	_, err := r.db.ExecContext(ctx, query, promotionID)
	if err != nil {
		return err
	}

	return nil
}

func NewPromotionRepository(db *sqlx.DB) PromotionRepository {
	return &promotionRepository{
		db: db,
	}
}
