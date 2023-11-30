package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/model"
)

type (
	ProductMediaRepository interface {
		FindProductMediaByProductID(ctx context.Context, id int64) ([]model.ProductMedia, error)
	}
	productMediaRepository struct {
		db *sqlx.DB
	}
)

func (r *productMediaRepository) FindProductMediaByProductID(ctx context.Context, id int64) ([]model.ProductMedia, error) {
	productMedias := make([]model.ProductMedia, 0)
	query := `SELECT * FROM product_medias pm WHERE pm.product_id = $1`
	err := r.db.SelectContext(ctx, &productMedias, query, id)
	if err != nil {
		return nil, err
	}
	return productMedias, nil
}

func NewProductMediaRepository(db *sqlx.DB) ProductMediaRepository {
	return &productMediaRepository{
		db: db,
	}
}
