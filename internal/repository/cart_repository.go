package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
)

type (
	CartRepository interface {
		FindCheckedCartByAccountID(ctx context.Context, accountID int64) ([]dto.CartOrderModel, error)
		FirstCart(ctx context.Context, cartId int64) (*model.Cart, error)
		FindByAccountID(ctx context.Context, accountId int64) ([]dto.CartPageModel, error)
		CountCartByAccountID(ctx context.Context, accountID int64) (*int64, error)
		FirstByProductVariantID(ctx context.Context, pVariantId, accountID int64) (*model.Cart, error)
		IncreaseQuantityByID(ctx context.Context, cartId int64, amount int) error
		Create(ctx context.Context, item *model.Cart) error
		UpdateQuantity(ctx context.Context, quantity int, cartId int64) error
		DeleteCart(ctx context.Context, cartId int64) error
		UpdateCheck(ctx context.Context, items []model.Cart) error
		FindCheckedForPrice(ctx context.Context, accountId int64) ([]dto.IsCheckedModel, error)
		FindCheckedCartByShopID(ctx context.Context, accountId, shopId int64) ([]dto.CartOrderModel, error)
		FindCartForHome(ctx context.Context, accountId int64) ([]dto.CartHomeModel, error)
	}
	cartRepository struct {
		db *sqlx.DB
	}
)

func (r *cartRepository) FindCheckedCartByAccountID(ctx context.Context, accountID int64) ([]dto.CartOrderModel, error) {
	qs := `
	SELECT 
			c.seller_id,
			s.id AS shop_id,
			s.name AS shop_name,
			c.id AS cart_id,
			p."name" AS product_name,
			p.id AS product_id,
			p.thumbnail_url AS image_url,
			c.product_variant_id,
			pv.price AS base_price,
			pv.discount,
			c.quantity,
			pv.stock AS remaining_quantity,
			vt."name" AS variant1_name,
			vt2."name" AS variant2_name,
			c.is_checked
		FROM carts c
		LEFT JOIN accounts a ON a.id = c.seller_id 
		LEFT JOIN shops s ON a.id = s.account_id 
		LEFT JOIN product_variants pv ON pv.id = c.product_variant_id
		LEFT JOIN products p ON p.id = pv.product_id
		LEFT JOIN variant_types vt ON vt.id = pv.variant_type1_id
		LEFT JOIN variant_types vt2 ON vt2.id = pv.variant_type2_id 
		WHERE c.account_id = $1 AND c.is_checked = TRUE
		ORDER BY s.name
	`

	com := make([]dto.CartOrderModel, 0)
	err := r.db.SelectContext(ctx, &com, qs, accountID)
	if err != nil {
		return nil, err
	}
	return com, nil
}

func (r *cartRepository) FirstCart(ctx context.Context, cartId int64) (*model.Cart, error) {
	cartProduct := new(model.Cart)
	query := `SELECT * FROM carts c WHERE c.id = $1 AND c.deleted_at is NULL`
	err := r.db.GetContext(ctx, cartProduct, query, cartId)
	if err != nil {
		return nil, err
	}
	return cartProduct, nil
}

// CountCartByAccountID implements CartRepository.
func (r *cartRepository) CountCartByAccountID(ctx context.Context, accountID int64) (*int64, error) {
	qs := `
		SELECT 
			COUNT(1) as cart_count 
		FROM carts c
		LEFT JOIN product_variants pv
			ON pv.product_id = c.id
		LEFT JOIN accounts buyer
			ON buyer.id = c.account_id
		LEFT JOIN accounts seller
			ON seller.id = c.seller_id
		WHERE c.id = $1;
	`
	row := r.db.QueryRowxContext(ctx, qs, accountID)
	if err := row.Err(); err != nil {
		return nil, err
	}
	count := new(int64)
	row.Scan(count)

	return count, nil
}

func (r *cartRepository) FindCheckedForPrice(ctx context.Context, accountId int64) ([]dto.IsCheckedModel, error) {
	carts := make([]dto.IsCheckedModel, 0)
	query := `
	SELECT 
		c.quantity,
		pv.price,
		pv.discount 
	FROM carts c 
	LEFT JOIN product_variants pv ON c.product_variant_id = pv.id 
	WHERE c.is_checked = TRUE AND c.account_id = $1
`
	err := r.db.SelectContext(ctx, &carts, query, accountId)
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (r *cartRepository) FindByAccountID(ctx context.Context, accountId int64) ([]dto.CartPageModel, error) {
	cartProducts := make([]dto.CartPageModel, 0)
	query := `
		SELECT 
			s.name AS shop_name,
			s.id AS shop_id,
			c.id AS cart_id,
			p."name" AS product_name,
			p.id AS product_id,
			p.thumbnail_url AS image_url,
			pv.price AS base_price,
			pv.discount,
			c.quantity,
			pv.stock AS remaining_quantity,
			vt."name" AS variant1_name,
			vt2."name" AS variant2_name,
			c.is_checked
		FROM carts c 
		LEFT JOIN accounts a ON a.id = c.seller_id 
		LEFT JOIN shops s ON a.id = s.account_id 
		LEFT JOIN product_variants pv ON pv.id = c.product_variant_id
		LEFT JOIN products p ON p.id = pv.product_id
		LEFT JOIN variant_types vt ON vt.id = pv.variant_type1_id
		LEFT JOIN variant_types vt2 ON vt2.id = pv.variant_type2_id 
		WHERE c.account_id = $1
		ORDER BY s.name, c.id
	`
	err := r.db.SelectContext(ctx, &cartProducts, query, accountId)
	if err != nil {
		return nil, err
	}
	return cartProducts, nil
}

func (r *cartRepository) FirstByProductVariantID(ctx context.Context, pVariantId, accountID int64) (*model.Cart, error) {
	cartProduct := new(model.Cart)
	query := `SELECT * FROM carts c WHERE c.product_variant_id = $1 AND c.account_id = $2`
	err := r.db.GetContext(ctx, cartProduct, query, pVariantId, accountID)
	if err != nil {
		return nil, err
	}
	return cartProduct, nil
}

func (r *cartRepository) Create(ctx context.Context, item *model.Cart) error {
	query := `
		INSERT INTO carts (
			quantity, 
			product_variant_id, 
			account_id,
			seller_id
			) VALUES (
				$1,
				$2,
				$3,
				$4
			)
	`
	_, err := r.db.ExecContext(ctx, query, item.Quantity, item.ProductVariantId, item.AccountId, item.SellerId)
	if err != nil {
		return err
	}

	return nil
}

func (r *cartRepository) IncreaseQuantityByID(ctx context.Context, cartId int64, amount int) error {
	query := `
		UPDATE carts
		SET quantity = quantity + $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, amount, time.Now(), cartId)
	if err != nil {
		return err
	}

	return nil
}

func (r *cartRepository) UpdateQuantity(ctx context.Context, quantity int, cartId int64) error {
	query := `
		UPDATE carts
		SET quantity = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, quantity, time.Now(), cartId)
	if err != nil {
		return err
	}

	return nil
}

func (r *cartRepository) UpdateCheck(ctx context.Context, items []model.Cart) error {
	query := `
		UPDATE carts
		SET is_checked = $1, updated_at = $2
		WHERE id = $3
	`
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	for _, item := range items {
		_, err := r.db.ExecContext(ctx, query, item.IsChecked, time.Now(), item.ID)
		if err != nil {
			err = tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *cartRepository) DeleteCart(ctx context.Context, cartId int64) error {
	query := `
		DELETE FROM carts
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, cartId)
	if err != nil {
		return err
	}

	return nil
}

func (r *cartRepository) FindCheckedCartByShopID(ctx context.Context, accountId, shopId int64) ([]dto.CartOrderModel, error) {
	cartOrders := make([]dto.CartOrderModel, 0)
	query := `
		SELECT 
			c.seller_id,
			s.id AS shop_id,
			s.name AS shop_name,
			c.id AS cart_id,
			p."name" AS product_name,
			p.id AS product_id,
			p.thumbnail_url AS image_url,
			c.product_variant_id,
			pv.price AS base_price,
			pv.discount,
			c.quantity,
			pv.stock AS remaining_quantity,
			vt."name" AS variant1_name,
			vt2."name" AS variant2_name,
			c.is_checked
		FROM carts c 
		LEFT JOIN accounts a ON a.id = c.seller_id 
		LEFT JOIN shops s ON a.id = s.account_id 
		LEFT JOIN product_variants pv ON pv.id = c.product_variant_id
		LEFT JOIN products p ON p.id = pv.product_id
		LEFT JOIN variant_types vt ON vt.id = pv.variant_type1_id
		LEFT JOIN variant_types vt2 ON vt2.id = pv.variant_type2_id 
		WHERE c.account_id = $1 AND c.is_checked = $2 AND s.id = $3
		ORDER BY a.id
	`
	err := r.db.SelectContext(ctx, &cartOrders, query, accountId, true, shopId)
	if err != nil {
		return nil, err
	}
	return cartOrders, nil
}

func (r *cartRepository) FindCartForHome(ctx context.Context, accountId int64) ([]dto.CartHomeModel, error) {
	carts := make([]dto.CartHomeModel, 0)
	query := `
		SELECT 
			p."name" AS product_name,
			p.thumbnail_url,
			pv.price AS base_price,
			pv.discount,
			c.quantity
		FROM carts c 
		LEFT JOIN product_variants pv ON pv.id = c.product_variant_id
		LEFT JOIN products p ON p.id = pv.product_id
		WHERE c.account_id = $1
		ORDER BY c.created_at
		LIMIT 5
	`
	err := r.db.SelectContext(ctx, &carts, query, accountId)
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func NewCartRepository(db *sqlx.DB) CartRepository {
	return &cartRepository{
		db: db,
	}
}
