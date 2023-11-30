package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
)

type (
	VariantGroupRepository interface {
		FindVariantGroupByProductID(ctx context.Context, id int64) ([]model.VariantGroup, error)
		FindVariantGroupWIthVariantTypeByProductID(ctx context.Context, id int64) ([]dto.ProductPageVariants, error)
	}
	variantGroupRepository struct {
		db *sqlx.DB
	}
)

func (r *variantGroupRepository) FindVariantGroupByProductID(ctx context.Context, id int64) ([]model.VariantGroup, error) {
	variantGroups := make([]model.VariantGroup, 0)
	query := `SELECT * FROM variant_groups vg WHERE vg.product_id = $1`
	err := r.db.SelectContext(ctx, &variantGroups, query, id)
	if err != nil {
		return nil, err
	}
	return variantGroups, nil
}

func (r *variantGroupRepository) FindVariantGroupWIthVariantTypeByProductID(ctx context.Context, id int64) ([]dto.ProductPageVariants, error) {
	variantGroups := make([]dto.ProductPageVariants, 0)
	query := `
		SELECT 
			vg.name AS group_name, 
			vt.id type_id, 
			vt.name AS type_name 
		FROM variant_groups vg 
		LEFT JOIN variant_types vt ON vt.variant_group_id = vg.id 
		WHERE vg.product_id = $1
		ORDER BY vg.id, CASE WHEN vt."name" = 'default' then 1 END
		`
	err := r.db.SelectContext(ctx, &variantGroups, query, id)
	if err != nil {
		return nil, err
	}
	return variantGroups, nil
}

func NewVariantGroupRepository(db *sqlx.DB) VariantGroupRepository {
	return &variantGroupRepository{
		db: db,
	}
}
