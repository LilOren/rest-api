package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type (
	OrderDetailRepository interface {
		CountOrderByProductCode(ctx context.Context, productCode string) (*int, error)
	}
	orderDetailRepository struct {
		db *sqlx.DB
	}
)

// CountOrderByProductCode implements OrderDetailRepository.
func (r *orderDetailRepository) CountOrderByProductCode(ctx context.Context, productCode string) (*int, error) {

	qs := `
	SELECT
		COUNT(1)
	FROM
		order_details od
	WHERE
		od.product_code = $1
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

func NewOrderDetailRepository(db *sqlx.DB) OrderDetailRepository {
	return &orderDetailRepository{
		db: db,
	}
}
