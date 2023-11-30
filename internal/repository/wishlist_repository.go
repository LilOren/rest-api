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
	WishlistRepository interface {
		FirstByUserAndProduct(ctx context.Context, userID, productID int64) (*model.Wishlist, error)
		FirstByID(ctx context.Context, wishlistID int64) (*model.Wishlist, error)
		FindByAccountID(ctx context.Context, userID int64, params *dto.WishlistParams) ([]dto.WishlistUserModel, error)
		FindByAccountIDMetadata(ctx context.Context, userID int64) ([]model.Wishlist, error)
		Create(ctx context.Context, wishlist *model.Wishlist) error
		Delete(ctx context.Context, wishlistID int64) error
		CountByProductID(ctx context.Context, prodID int64) (*dto.WishlistCountPayload, error)
	}
	wishlistRepository struct {
		db *sqlx.DB
	}
)

func (r *wishlistRepository) FirstByUserAndProduct(ctx context.Context, userID, productID int64) (*model.Wishlist, error) {
	wishlist := new(model.Wishlist)
	query := `SELECT * FROM wishlists w WHERE w.account_id = $1 AND w.product_id = $2`
	err := r.db.GetContext(ctx, wishlist, query, userID, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return wishlist, nil
}

func (r *wishlistRepository) FirstByID(ctx context.Context, wishlistID int64) (*model.Wishlist, error) {
	wishlist := new(model.Wishlist)
	query := `SELECT * FROM wishlists w WHERE w.id = $1`
	err := r.db.GetContext(ctx, wishlist, query, wishlistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return wishlist, nil
}

func (r *wishlistRepository) FindByAccountID(ctx context.Context, userID int64, params *dto.WishlistParams) ([]dto.WishlistUserModel, error) {
	wishlist := make([]dto.WishlistUserModel, 0)
	query := `
	SELECT
		w.id,
		p.product_code AS product_code,
		p.thumbnail_url as thumbnail_url,
		p.name AS product_name,
		pv.price as base_price,
		pv.discount as discount,
		s.name AS shop_name,
		d.name AS district_name
	FROM wishlists w 
	LEFT JOIN products p ON w.product_id = p.id
	LEFT JOIN (
		SELECT
			DISTINCT ON
			(pv.product_id) pv.product_id,
			pv.discount,
			pv.price
		FROM
			product_variants pv
		ORDER BY
			pv.product_id ASC
	) pv ON p.id = pv.product_id
	LEFT JOIN shops s ON s.account_id = p.seller_id
	LEFT JOIN account_addresses aa ON aa.account_id = p.seller_id AND aa.is_shop
	LEFT JOIN districts d ON aa.district_id = d.id
	WHERE w.account_id = $1
	ORDER BY w.created_at DESC
	LIMIT $2 OFFSET $3
	`
	limit := constant.WishlistDefaultItems
	offset := (params.Page - 1) * limit
	err := r.db.SelectContext(ctx, &wishlist, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return wishlist, nil
}

func (r *wishlistRepository) FindByAccountIDMetadata(ctx context.Context, userID int64) ([]model.Wishlist, error) {
	wishlist := make([]model.Wishlist, 0)
	query := `
	SELECT *
	FROM wishlists w 
	WHERE w.account_id = $1
	ORDER BY w.created_at DESC
	`
	err := r.db.SelectContext(ctx, &wishlist, query, userID)
	if err != nil {
		return nil, err
	}
	return wishlist, nil
}

func (r *wishlistRepository) Create(ctx context.Context, wishlist *model.Wishlist) error {
	query := `
	INSERT INTO wishlists (account_id, product_id) 
	VALUES 
	($1,$2)
	`
	_, err := r.db.ExecContext(ctx, query, wishlist.AccountID, wishlist.ProductID)
	if err != nil {
		return err
	}

	return nil
}

func (r *wishlistRepository) Delete(ctx context.Context, wishlistID int64) error {
	query := `
		DELETE FROM wishlists
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, wishlistID)
	if err != nil {
		return err
	}

	return nil
}

func (r *wishlistRepository) CountByProductID(ctx context.Context, prodID int64) (*dto.WishlistCountPayload, error) {
	wishlistCtr := new(dto.WishlistCountPayload)
	query := `
		SELECT 
			count(w.id) AS counter
		FROM wishlists w 
		WHERE w.product_id = $1
		GROUP BY w.product_id
	`
	err := r.db.GetContext(ctx, wishlistCtr, query, prodID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			wishlistCtr.Counter = 0
			return wishlistCtr, nil
		}
		return nil, err
	}
	return wishlistCtr, nil
}

func NewWishlistRepository(db *sqlx.DB) WishlistRepository {
	return &wishlistRepository{
		db: db,
	}
}
