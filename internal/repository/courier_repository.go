package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/model"
)

type (
	CourierRepository interface {
		FirstByShopCourierID(ctx context.Context, shopCourierID int64) (*model.Courier, error)
	}
	courierRepository struct {
		db *sqlx.DB
	}
)

// FirstByShopCourierID implements CourierRepository.
func (r *courierRepository) FirstByShopCourierID(ctx context.Context, shopCourierID int64) (*model.Courier, error) {
	courier := new(model.Courier)
	qs := `
	SELECT 
	c.* 
	FROM couriers c
	JOIN shop_couriers sc
		ON sc.courier_id = c.id
	WHERE sc.id = $1
	`
	err := r.db.GetContext(ctx, courier, qs, shopCourierID)
	if err != nil {
		return nil, err
	}

	return courier, nil
}

func NewCourierRepository(db *sqlx.DB) CourierRepository {
	return &courierRepository{
		db: db,
	}
}
