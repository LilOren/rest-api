package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
)

type (
	ShopRepository interface {
		FirstShopById(ctx context.Context, accountId int) (*model.Shop, error)
		FirstShopDetailByAccountID(ctx context.Context, id int64) (*dto.ProductPageShop, error)
		CreateShop(ctx context.Context, payload dto.CreateShopPayload, accountId int) error
		UpdateShopName(ctx context.Context, shopName string, accountId int) error
		UpdateShopAddress(ctx context.Context, addressId int, accountId int) error
		UpdateShopCourier(ctx context.Context, payload []bool, shopId int) error
		GetAllProductBySellerId(ctx context.Context, sellerId int, page int) ([]dto.GetAllProduct, error)
		CountAllProductBySellerId(ctx context.Context, sellerId int) (*int, error)
		FindAllProductDiscountByProductCode(ctx context.Context, sellerId int64, productCode string) ([]dto.GetProductDiscount, error)
		UpdateProductDiscount(ctx context.Context, payload dto.UpdateProductDiscountPayload, sellerId int64, productCode string) error
	}
	shopRepository struct {
		db *sqlx.DB
	}
)

func (r *shopRepository) FirstShopDetailByAccountID(ctx context.Context, id int64) (*dto.ProductPageShop, error) {
	shop := new(dto.ProductPageShop)
	query := `
	SELECT
		a.id, 
		s.name, 
		COALESCE(
			a.profile_picture_url,
			''
		) AS profile_picture_url, 
		d.name AS LOCATION
	FROM
		shops s
	LEFT JOIN accounts a ON
		s.account_id = a.id
	LEFT JOIN account_addresses ad ON
		ad.account_id = a.id
	LEFT JOIN districts d ON
		d.id = ad.district_id
	WHERE
		s.account_id = $1
		AND ad.is_shop
	`
	err := r.db.GetContext(ctx, shop, query, id)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (sr *shopRepository) FirstShopById(ctx context.Context, accountId int) (*model.Shop, error) {
	shop := new(model.Shop)
	qs := `
	SELECT 
		id,
		name,
		account_id
	FROM
		shops
	WHERE
		account_id = $1
	`

	row := sr.db.QueryRowxContext(ctx, qs, accountId)
	if err := row.Err(); err != nil {
		return nil, err
	}

	err := row.StructScan(shop)
	if err != nil {
		return nil, err
	}

	return shop, nil
}

func (sr *shopRepository) CreateShop(ctx context.Context, payload dto.CreateShopPayload, accountId int) error {
	qs1 := `
	UPDATE account_addresses
	SET is_shop = $1, updated_at = $2
	WHERE account_id = $3 AND id = $4
	`

	qs2 := `
	UPDATE accounts
	SET is_seller = $1, updated_at = $2
	WHERE id = $3
	`

	qs3 := `
	INSERT INTO shop_couriers (shop_id, courier_id, is_available) VALUES 
	($1,$2,$3),
	($4,$5,$6),
	($7,$8,$9)
	`

	tx, err := sr.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var id int64
	err = tx.QueryRow("INSERT INTO shops (name, account_id,updated_at) VALUES ($1,$2,$3) RETURNING id", payload.ShopName, accountId, time.Now()).Scan(&id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(qs1, true, time.Now(), accountId, payload.AddressId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(qs2, true, time.Now(), accountId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(qs3, id, 1, true, id, 2, false, id, 3, false)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (sr *shopRepository) UpdateShopName(ctx context.Context, shopName string, accountId int) error {
	qs := `
	UPDATE shops
	SET name = $1, updated_at = $2
	WHERE account_id = $3
	`

	_, err := sr.db.ExecContext(ctx, qs, shopName, time.Now(), accountId)
	if err != nil {
		return err
	}

	return nil
}

func (sr *shopRepository) UpdateShopAddress(ctx context.Context, addressId int, accountId int) error {
	qs1 := `
	UPDATE account_addresses
	SET is_shop = $1, updated_at = $2
	WHERE account_id = $3 AND id = $4
	`

	qs2 := `
	UPDATE account_addresses
	SET is_shop = $1, updated_at = $2
	WHERE account_id = $3 AND id != $4
	`

	tx, err := sr.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(qs1, true, time.Now(), accountId, addressId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(qs2, false, time.Now(), accountId, addressId)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (sr *shopRepository) UpdateShopCourier(ctx context.Context, payload []bool, shopId int) error {
	qs := `
	UPDATE shop_couriers
	SET is_available = $1, updated_at = $2
	WHERE shop_id = $3 AND courier_id = $4
	`

	tx, err := sr.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for k, v := range payload {
		_, err = tx.ExecContext(ctx, qs, v, time.Now(), shopId, k+1)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *shopRepository) GetAllProductBySellerId(ctx context.Context, sellerId int, page int) ([]dto.GetAllProduct, error) {
	products := make([]dto.GetAllProduct, 0)
	query := `
		SELECT 
			p.product_code, 
			p.name AS product_name, 
			p.thumbnail_url
		FROM
			products p
		WHERE
			p.seller_id = $1
		ORDER BY
			p.updated_at DESC
		LIMIT 10
		OFFSET $2
		`

	rows, err := r.db.QueryxContext(ctx, query, sellerId, (page-1)*10)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := new(dto.GetAllProduct)
		if err := rows.StructScan(temp); err != nil {
			return nil, err
		}
		products = append(products, *temp)
	}

	return products, nil
}

func (r *shopRepository) CountAllProductBySellerId(ctx context.Context, sellerId int) (*int, error) {
	products := make([]dto.GetAllProduct, 0)
	query := `
		SELECT 
			p.product_code, 
			p.name AS product_name, 
			p.thumbnail_url
		FROM
			products p
		WHERE
			p.seller_id = $1
		ORDER BY
			p.created_at DESC
	`

	rows, err := r.db.QueryxContext(ctx, query, sellerId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := new(dto.GetAllProduct)
		if err := rows.StructScan(temp); err != nil {
			return nil, err
		}
		products = append(products, *temp)
	}

	res := len(products)

	return &res, nil
}

func (r *shopRepository) FindAllProductDiscountByProductCode(ctx context.Context, sellerId int64, productCode string) ([]dto.GetProductDiscount, error) {
	discounts := make([]dto.GetProductDiscount, 0)
	query := `
	SELECT
		pv.discount,
		vg1."name" AS variant_group1_name,
		vt1."name" AS variant_type1_name,
		vg2."name" AS variant_group2_name,
		vt2."name" AS variant_type2_name 
	FROM
		product_variants pv 
		LEFT JOIN products p ON pv.product_id = p.id 
		LEFT JOIN variant_types vt1 ON pv.variant_type1_id = vt1.id 
		LEFT JOIN variant_groups vg1 ON vg1.id = vt1.variant_group_id  
		LEFT JOIN variant_types vt2 ON pv.variant_type2_id =  vt2.id 
		LEFT JOIN variant_groups vg2 ON vg2.id = vt2.variant_group_id
	WHERE
		p.seller_id = $1 AND
		p.product_code = $2
	ORDER BY
		pv.id ASC;
	`

	rows, err := r.db.QueryxContext(ctx, query, sellerId, productCode)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := new(dto.GetProductDiscount)
		if err := rows.StructScan(temp); err != nil {
			return nil, err
		}
		discounts = append(discounts, *temp)
	}

	return discounts, nil
}

func (r *shopRepository) UpdateProductDiscount(ctx context.Context, payload dto.UpdateProductDiscountPayload, sellerId int64, productCode string) error {
	query := `
	UPDATE product_variants
	SET discount = $1, updated_at = now()
	WHERE id = (
		SELECT
			pv.id
		FROM
			product_variants pv 
			LEFT JOIN products p ON pv.product_id = p.id 
			LEFT JOIN variant_types vt1 ON pv.variant_type1_id = vt1.id 
			LEFT JOIN variant_groups vg1 ON vg1.id = vt1.variant_group_id  
			LEFT JOIN variant_types vt2 ON pv.variant_type2_id =  vt2.id 
			LEFT JOIN variant_groups vg2 ON vg2.id = vt2.variant_group_id
		WHERE
			p.seller_id = $2 AND
			vt1.name = $3 AND
			vt2.name = $4 AND
			p.product_code = $5
	)
	`
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, v := range payload.Variants {
		_, err := tx.Exec(query, v.Discount, sellerId, v.VariantType1, v.VariantType2, productCode)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func NewShopRepository(db *sqlx.DB) ShopRepository {
	return &shopRepository{
		db: db,
	}
}
