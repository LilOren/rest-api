package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/model"
)

type (
	ProvinceRepository interface {
		Find(ctx context.Context) ([]model.Province, error)
		FirstByID(ctx context.Context, provinceID int64) (*model.Province, error)
	}
	provinceRepository struct {
		db *sqlx.DB
	}
)

// FirstByID implements ProvinceRepository.
func (r *provinceRepository) FirstByID(ctx context.Context, provinceID int64) (*model.Province, error) {
	qs := `
	SELECT
		p.*
	FROM 
		provinces p
	WHERE
		p.id = $1
	`

	row := r.db.QueryRowxContext(ctx, qs, provinceID)

	if err := row.Err(); err != nil {
		return nil, err
	}

	m := new(model.Province)
	if err := row.StructScan(m); err != nil {
		return nil, err
	}

	return m, nil
}

// Find implements ProvinceRepository.
func (r *provinceRepository) Find(ctx context.Context) ([]model.Province, error) {
	province := make([]model.Province, 0)

	qs := `SELECT * FROM provinces`

	err := r.db.SelectContext(ctx, &province, qs)
	if err != nil {
		return nil, err
	}

	return province, nil
}

func NewProvinceRepository(db *sqlx.DB) ProvinceRepository {
	return &provinceRepository{
		db: db,
	}
}
