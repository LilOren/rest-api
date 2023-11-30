package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/model"
)

type (
	DistrictRepository interface {
		FindByProvinceID(ctx context.Context, provinceId int64) ([]model.District, error)
		FirstByID(ctx context.Context, districtID int64) (*model.District, error)
	}
	districtRepository struct {
		db *sqlx.DB
	}
)

// FirstByID implements DistrictRepository.
func (r *districtRepository) FirstByID(ctx context.Context, districtID int64) (*model.District, error) {
	qs := `
	SELECT
		d.*
	FROM districts d
	WHERE
		d.id = $1
	`

	row := r.db.QueryRowxContext(ctx, qs, districtID)
	if err := row.Err(); err != nil {
		return nil, err
	}

	m := new(model.District)
	if err := row.StructScan(m); err != nil {
		return nil, err
	}

	return m, nil
}

// Find implements ProvinceRepository.
func (r *districtRepository) FindByProvinceID(ctx context.Context, provinceId int64) ([]model.District, error) {
	districts := make([]model.District, 0)

	qs := `SELECT * FROM districts WHERE province_id = $1`

	err := r.db.SelectContext(ctx, &districts, qs, provinceId)
	if err != nil {
		return nil, err
	}

	return districts, nil
}

func NewDistrictRepository(db *sqlx.DB) DistrictRepository {
	return &districtRepository{
		db: db,
	}
}
