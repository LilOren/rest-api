package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
)

type (
	ProductVariantRepository interface {
		FindProductVariantByProductID(ctx context.Context, accountId int64) ([]model.ProductVariant, error)
		FirstProductVariantByIDForCart(ctx context.Context, id int64) (*dto.ProductVariantForCartModel, error)
		FirstProductVariantByID(ctx context.Context, id int64) (*model.ProductVariant, error)
	}
	productVariantRepository struct {
		db *sqlx.DB
	}
)

func (r *productVariantRepository) FindProductVariantByProductID(ctx context.Context, productID int64) ([]model.ProductVariant, error) {
	productVariant := make([]model.ProductVariant, 0)
	query := `SELECT * FROM product_variants pv WHERE pv.product_id = $1`
	err := r.db.SelectContext(ctx, &productVariant, query, productID)
	if err != nil {
		return nil, err
	}
	return productVariant, nil
}

func (r *productVariantRepository) FirstProductVariantByIDForCart(ctx context.Context, id int64) (*dto.ProductVariantForCartModel, error) {
	productVariant := new(dto.ProductVariantForCartModel)
	query := `
		SELECT 
			pv.id,
			pv.stock, 
			p.seller_id
		FROM product_variants pv 
		LEFT JOIN products p ON p.id  = pv.product_id 
		WHERE pv.id = $1
	`
	err := r.db.GetContext(ctx, productVariant, query, id)
	if err != nil {
		return nil, err
	}
	return productVariant, nil
}

func (r *productVariantRepository) FirstProductVariantByID(ctx context.Context, id int64) (*model.ProductVariant, error) {
	productVariant := new(model.ProductVariant)
	query := `SELECT * FROM product_variants pv WHERE pv.id = $1`
	err := r.db.GetContext(ctx, productVariant, query, id)
	if err != nil {
		return nil, err
	}
	return productVariant, nil
}

func DecreaseStock(tx *sqlx.Tx, code string, quantity int, vName1, vName2 string) error {
	qs := `
	UPDATE product_variants
	SET stock = stock-$1, updated_at = now()
	WHERE id = ( 
			SELECT
				pv.id  
			FROM product_variants pv 
			LEFT JOIN products p ON p.id = pv.product_id 
			LEFT JOIN variant_types vt ON pv.variant_type1_id = vt.id
			LEFT JOIN variant_types vt2 ON pv.variant_type2_id = vt2.id 
			WHERE p.product_code = $2
			AND (vt."name" = $3 AND vt2."name" = $4)
	)
	`
	_, err := tx.Exec(qs, quantity, code, vName1, vName2)
	if err != nil {
		return err
	}

	return nil
}

func NewProductVariantRepository(db *sqlx.DB) ProductVariantRepository {
	return &productVariantRepository{
		db: db,
	}
}
