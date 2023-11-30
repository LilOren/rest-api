package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/dto"
)

type (
	ShopCourierRepository interface {
		FindShopCourierByShopId(ctx context.Context, shopId int64) ([]dto.ShopCourierDetailsResponse, error)
		FindAvailableCourierByShopID(ctx context.Context, shopId int64) ([]dto.ShopCourierDetailsResponse, error)
	}
	shopCourierRepository struct {
		db *sqlx.DB
	}
)

// FindAvailableCourierByShopID implements ShopCourierRepository.
func (r *shopCourierRepository) FindAvailableCourierByShopID(ctx context.Context, shopId int64) ([]dto.ShopCourierDetailsResponse, error) {
	shopCourierList := make([]dto.ShopCourierDetailsResponse, 0)
	qs := `
		SELECT 
			sc.id as shop_courier_id,
			c.name as courier_name,
			sc.is_available
		FROM
			shop_couriers sc
			LEFT JOIN couriers c ON sc.courier_id = c.id
		WHERE
			sc.shop_id = $1 AND sc.is_available = TRUE
	`

	rows, err := r.db.QueryxContext(ctx, qs, shopId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := new(dto.ShopCourierDetailsResponse)
		rows.StructScan(temp)
		shopCourierList = append(shopCourierList, *temp)
	}

	return shopCourierList, nil
}

func (r *shopCourierRepository) FindShopCourierByShopId(ctx context.Context, shopId int64) ([]dto.ShopCourierDetailsResponse, error) {
	shopCourierList := make([]dto.ShopCourierDetailsResponse, 0)
	qs := `
		SELECT 
			sc.id as shop_courier_id,
			c.name as courier_name,
			c.image_url,
			c.description,
			sc.is_available
		FROM
			shop_couriers sc
			LEFT JOIN couriers c ON sc.courier_id = c.id
		WHERE
			shop_id = $1
	`

	rows, err := r.db.QueryxContext(ctx, qs, shopId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := new(dto.ShopCourierDetailsResponse)
		rows.StructScan(temp)
		shopCourierList = append(shopCourierList, *temp)
	}

	return shopCourierList, nil
}

func NewShopCourierRepository(db *sqlx.DB) ShopCourierRepository {
	return &shopCourierRepository{
		db: db,
	}
}
