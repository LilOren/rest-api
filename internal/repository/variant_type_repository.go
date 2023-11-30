package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/model"
)

type (
	VariantTypeRepository interface {
		FindVariantTypeByVariantGroupID(ctx context.Context, ids []int64) ([]model.VariantType, error)
	}
	variantTypeRepository struct {
		db *sqlx.DB
	}
)

func (r *variantTypeRepository) FindVariantTypeByVariantGroupID(ctx context.Context, ids []int64) ([]model.VariantType, error) {
	variantTypes := make([]model.VariantType, 0)
	query, args, err := sqlx.In(`SELECT * FROM variant_types vt WHERE vt.variant_group_id IN (?);`, ids)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		variantType := new(model.VariantType)
		rows.StructScan(variantType)
		variantTypes = append(variantTypes, *variantType)
	}
	return variantTypes, nil
}

func NewVariantTypeRepository(db *sqlx.DB) VariantTypeRepository {
	return &variantTypeRepository{
		db: db,
	}
}
