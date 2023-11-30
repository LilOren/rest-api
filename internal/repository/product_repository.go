package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
	"github.com/lil-oren/rest/internal/shared"
)

type (
	ProductRepository interface {
		FindRecommended(ctx context.Context) ([]dto.HomePageProductModel, error)
		FirstProductDetail(ctx context.Context, id int64) (*model.Product, error)
		FindProductBySearchTerm(ctx context.Context, payload dto.SearchProductPayload) ([]dto.SearchProductResponseItem, error)
		CountProductBySearchTerm(ctx context.Context, payload dto.CountProductBySearchTermPayload) (*int, error)
		CreateProduct(ctx context.Context, payload dto.AddProduct, accountId int, productCode string, mediaType []string) error
		FindProductDetail(ctx context.Context, productCode string) ([]dto.GetProductDetail, error)
		FindProductVariants(ctx context.Context, productId int) ([]dto.GetProductVariant, error)
		FirstProductByCode(ctx context.Context, code string) (*model.Product, error)
		UpdateProduct(ctx context.Context, payload dto.UpdateProductPayload, categories []interface{}, mediaType []string, accountId int64) error
		DeleteProduct(ctx context.Context, productCode string, sellerId int64) error
	}
	productRepository struct {
		db *sqlx.DB
	}
)

func (r *productRepository) DeleteProduct(ctx context.Context, productCode string, sellerId int64) error {
	query := `
	DELETE FROM products
	WHERE id = (
		SELECT
			p.id
		FROM
			products p
		WHERE
			p.product_code = $1 AND p.seller_id = $2
	)
	`

	_, err := r.db.ExecContext(ctx, query, productCode, sellerId)
	if err != nil {
		return err
	}

	return nil
}

// CountProductBySearchTerm implements ProductRepository.
func (r *productRepository) CountProductBySearchTerm(ctx context.Context, payload dto.CountProductBySearchTermPayload) (*int, error) {
	searchTerm := "%"
	if payload.SearchTerm != "" {
		searchTerm = "%" + payload.SearchTerm + "%"
	}

	qs := `
	SELECT
		COUNT(1) as count
	FROM
		(
			SELECT 
				DISTINCT ON (p.id)
				p.*,
				pc.category_id
			FROM
				products p
			LEFT JOIN product_categories pc
				ON pc.product_id = p.id
			WHERE 
				p.name ILIKE :search_term
				%s
		) p
	LEFT JOIN
			(
			SELECT
				DISTINCT ON
				(pv.product_id) pv.product_id,
				pv.discount,
				pv.price
			FROM
				product_variants pv
			ORDER BY
				pv.product_id ASC
		) pv ON
		p.id = pv.product_id
	LEFT JOIN shops s
		ON
		s.account_id = p.seller_id
	LEFT JOIN account_addresses aa
		ON
		s.account_id = aa.account_id
		AND aa.is_shop
	LEFT JOIN districts d
		ON
		aa.district_id = d.id
	WHERE
		%s
		AND
		%s
		AND
		%s
	`

	args := map[string]interface{}{
		"search_term": searchTerm,
	}

	if payload.CategoryID > 0 {
		condition := "AND pc.category_id = %d"
		condition = fmt.Sprintf(condition, payload.CategoryID)
		qs = fmt.Sprintf(qs, condition, "%s", "%s", "%s")
	} else {
		qs = fmt.Sprintf(qs, "", "%s", "%s", "%s")
	}

	if payload.DistrictIDs != "" {
		condition := "aa.district_id IN (%s)"
		condition = fmt.Sprintf(condition, payload.DistrictIDs)
		qs = fmt.Sprintf(qs, condition, "%s", "%s")
	} else {
		qs = fmt.Sprintf(qs, "TRUE", "%s", "%s")
	}

	if payload.MinPrice > 0 {
		condition := "pv.price >= %f"
		condition = fmt.Sprintf(condition, payload.MinPrice)
		qs = fmt.Sprintf(qs, condition, "%s")
	} else {
		qs = fmt.Sprintf(qs, "TRUE", "%s")
	}

	if payload.MaxPrice > 0 {
		condition := "pv.price <= %f"
		condition = fmt.Sprintf(condition, payload.MaxPrice)
		qs = fmt.Sprintf(qs, condition)
	} else {
		qs = fmt.Sprintf(qs, "TRUE")
	}

	rows, err := r.db.NamedQueryContext(ctx, qs, args)
	if err != nil {
		return nil, err
	}

	rows.Next()

	count := new(int)
	if err := rows.Scan(count); err != nil {
		return nil, err
	}

	return count, nil
}

// FindProductBySearchTerm implements ProductRepository.
func (r *productRepository) FindProductBySearchTerm(ctx context.Context, payload dto.SearchProductPayload) ([]dto.SearchProductResponseItem, error) {
	qs := `
	SELECT
		p.product_code AS product_code,
		p.thumbnail_url as thumbnail_url,
		p.name AS product_name,
		pv.price as base_price,
		pv.discount as discount,
		(pv.price - (pv.price * pv.discount / 100)) AS discount_price,
		s.name AS shop_name,
		d.name AS district_name,
		(CASE 
			WHEN mp.count_purchased IS NULL THEN 0
			ELSE mp.count_purchased
		END) AS total_sold
	FROM
		(
			SELECT 
				DISTINCT ON (p.id)
				p.*,
				pc.category_id
			FROM
				products p
			LEFT JOIN product_categories pc
				ON pc.product_id = p.id
			WHERE 
				p.name ILIKE :search_term
				%s
		) p
	LEFT JOIN
			(
			SELECT
				DISTINCT ON
				(pv.product_id) pv.product_id,
				pv.discount,
				pv.price
			FROM
				product_variants pv
			ORDER BY
				pv.product_id ASC
		) pv ON
		p.id = pv.product_id
	LEFT JOIN 
	(
		SELECT
			od.product_code,
			count(od.order_id) AS count_purchased
		FROM order_details od 
		GROUP BY od.product_code
	) mp ON 
	p.product_code = mp.product_code
	LEFT JOIN shops s
		ON
		s.account_id = p.seller_id
	LEFT JOIN account_addresses aa
		ON
		s.account_id = aa.account_id
		AND aa.is_shop
	LEFT JOIN districts d
		ON
		aa.district_id = d.id
	WHERE
		%s
		AND
		%s
		AND
		%s
	ORDER BY
		%s
	OFFSET :offset
	LIMIT 30
	`

	start := 0
	searchTerm := "%"

	if payload.SearchTerm != "" {
		searchTerm = "%" + payload.SearchTerm + "%"
	}

	if payload.Page > 1 {
		start = (payload.Page - 1) * 30
	}

	args := map[string]interface{}{
		"offset":      start,
		"search_term": searchTerm,
	}

	if payload.CategoryID > 0 {
		condition := "AND pc.category_id = %d"
		condition = fmt.Sprintf(condition, payload.CategoryID)
		qs = fmt.Sprintf(qs, condition, "%s", "%s", "%s", "%s")
	} else {
		qs = fmt.Sprintf(qs, "", "%s", "%s", "%s", "%s")
	}

	if payload.DistrictIDs != "" {
		condition := "aa.district_id IN (%s)"
		condition = fmt.Sprintf(condition, payload.DistrictIDs)
		qs = fmt.Sprintf(qs, condition, "%s", "%s", "%s")
	} else {
		qs = fmt.Sprintf(qs, "TRUE", "%s", "%s", "%s")
	}

	if payload.MinPrice > 0 {
		condition := "pv.price >= %f"
		condition = fmt.Sprintf(condition, payload.MinPrice)
		qs = fmt.Sprintf(qs, condition, "%s", "%s")
	} else {
		qs = fmt.Sprintf(qs, "TRUE", "%s", "%s")
	}

	if payload.MaxPrice > 0 {
		condition := "pv.price <= %f"
		condition = fmt.Sprintf(condition, payload.MaxPrice)
		qs = fmt.Sprintf(qs, condition, "%s")
	} else {
		qs = fmt.Sprintf(qs, "TRUE", "%s")
	}

	switch payload.SortBy {
	case "created_at":
		{
			if payload.SortDesc {
				qs = fmt.Sprintf(qs, "p.created_at DESC")
				break
			}

			qs = fmt.Sprintf(qs, "p.created_at ASC")
		}
	case "price":
		{
			if payload.SortDesc {
				qs = fmt.Sprintf(qs, "base_price DESC")
				break
			}

			qs = fmt.Sprintf(qs, "base_price ASC")
		}
	case "most_purchased":
		{
			if payload.SortDesc {
				qs = fmt.Sprintf(qs, "total_sold DESC")
				break
			}

			qs = fmt.Sprintf(qs, "total_sold ASC")
		}
	}

	stmt, err := r.db.PrepareNamedContext(ctx, qs)
	if err != nil {
		return nil, err
	}

	m := make([]dto.SearchProductResponseItem, 0)
	if err := stmt.SelectContext(ctx, &m, args); err != nil {
		return nil, err
	}

	return m, nil
}

func (r *productRepository) FindRecommended(ctx context.Context) ([]dto.HomePageProductModel, error) {
	e := make([]dto.HomePageProductModel, 0)

	query := `
	SELECT
		p.product_code,
		pm.media_url,
		p.name, 
		pv.price, 
		(pv.price - (pv.price * pv.discount / 100)) as discounted_price,
		pv.discount,
		(CASE 
			WHEN mp.count_purchased IS NULL THEN 0
			ELSE mp.count_purchased
		END) AS total_sold,
		d.name AS shop_location,
		s.name AS shop_name
	FROM products p
	LEFT JOIN 
		(
			SELECT 
				DISTINCT ON (pv.product_id) pv.product_id, 
				pv.discount, 
				pv.price 
			FROM 
				product_variants pv
		) pv ON p.id = pv.product_id
	LEFT JOIN 
	(
		SELECT
			od.product_code,
			count(od.order_id) AS count_purchased
		FROM order_details od 
		GROUP BY od.product_code
	) mp ON 
	p.product_code = mp.product_code
	LEFT JOIN product_medias pm 
		ON p.id = pm.product_id
	LEFT JOIN shops s
		ON s.account_id = p.seller_id
	LEFT JOIN account_addresses aa
		ON aa.account_id = p.seller_id
		AND aa.is_shop
	LEFT JOIN districts d
		ON d.id = aa.district_id
	LIMIT 18;
	`

	err := r.db.SelectContext(ctx, &e, query)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *productRepository) FirstProductDetail(ctx context.Context, id int64) (*model.Product, error) {
	product := new(model.Product)
	err := r.db.GetContext(ctx, product, "SELECT * FROM products p WHERE p.id = $1", id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepository) CreateProduct(ctx context.Context, payload dto.AddProduct, accountId int, productCode string, mediaType []string) error {
	qs1 := `
	INSERT INTO products (name, product_code, description, thumbnail_url, seller_id, weight) VALUES 
	($1,$2,$3,$4,$5,$6)
	RETURNING (id)
	`

	qs2 := `
	INSERT INTO product_medias (media_url, media_type, product_id) VALUES 
	($1,$2,$3)
	`

	qs3 := `
	INSERT INTO product_categories (product_id, category_id) VALUES 
	($1,$2)
	`

	qs4 := `
	INSERT INTO variant_groups (name, product_id) VALUES
	($1, $2)
	RETURNING (id)
	`

	qs5 := `
	INSERT INTO variant_types (name, variant_group_id) VALUES
	($1, $2)
	RETURNING (id)
	`

	qs6 := `
	INSERT INTO product_variants (price, stock, discount, product_id, variant_type1_id, variant_type2_id) VALUES
	($1, $2, $3, $4, $5, $6)
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var productID int64
	err = tx.QueryRowx(qs1, payload.ProductName, productCode, payload.Description, payload.ImageURL[0], accountId, payload.Weight).Scan(&productID)
	if err != nil {
		return err
	}

	for idx, v := range payload.ImageURL {
		_, err = tx.Exec(qs2, v, mediaType[idx], productID)
		if err != nil {
			return err
		}
	}

	if payload.ProductCategoryID.Level1 != 0 {
		_, err = tx.Exec(qs3, productID, payload.ProductCategoryID.Level1)
		if err != nil {
			return err
		}
	}

	if payload.ProductCategoryID.Level2 != 0 {
		_, err = tx.Exec(qs3, productID, payload.ProductCategoryID.Level2)
		if err != nil {
			return err
		}
	}

	if payload.ProductCategoryID.Level3 != nil {
		_, err = tx.Exec(qs3, productID, payload.ProductCategoryID.Level3)
		if err != nil {
			return err
		}
	}

	var variantGroup1ID, variantType1ID int64
	listVariantType1ID := make(map[string]int64)
	err = tx.QueryRowx(qs4, payload.VariantDefinitions.VariantGroup1.Name, productID).Scan(&variantGroup1ID)
	if err != nil {
		return err
	}

	for _, v := range payload.VariantDefinitions.VariantGroup1.VariantTypes {
		err = tx.QueryRowx(qs5, v, variantGroup1ID).Scan(&variantType1ID)
		if err != nil {
			return err
		}
		listVariantType1ID[v] = variantType1ID
	}

	var variantGroup2ID, variantType2ID int64
	listVariantType2ID := make(map[string]int64)
	err = tx.QueryRowx(qs4, payload.VariantDefinitions.VariantGroup2.Name, productID).Scan(&variantGroup2ID)
	if err != nil {
		return err
	}

	for _, v := range payload.VariantDefinitions.VariantGroup2.VariantTypes {
		err = tx.QueryRowx(qs5, v, variantGroup2ID).Scan(&variantType2ID)
		if err != nil {
			return err
		}
		listVariantType2ID[v] = variantType2ID
	}

	for _, v := range payload.Variants {
		_, err = tx.Exec(qs6, v.Price, v.Stock, 0, productID, listVariantType1ID[v.VariantType1], listVariantType2ID[v.VariantType2])
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *productRepository) FindProductDetail(ctx context.Context, productCode string) ([]dto.GetProductDetail, error) {
	detail := make([]dto.GetProductDetail, 0)
	qs := `
	SELECT
		p.id,
		p.product_code,
		p."name" AS product_name,
		p.description,
		p.weight,
		pm.id AS media_id,
		pm.media_url,
		pc.category_id,
		c."name" AS category_name,
		vt1.variant_group_id AS variant_group1_id,
		vg1.name AS variant_group1_name,
		vt2.variant_group_id AS variant_group2_id,
		vg2.name AS variant_group2_name
	FROM
		product_variants pv 
		LEFT JOIN products p ON pv.product_id = p.id
		LEFT JOIN variant_types vt1 ON pv.variant_type1_id = vt1.id
		LEFT JOIN variant_types vt2 ON pv.variant_type2_id = vt2.id
		LEFT JOIN variant_groups vg1 ON vg1.id = vt1.variant_group_id
		LEFT JOIN variant_groups vg2 ON vg2.id = vt2.variant_group_id
		LEFT JOIN product_medias pm ON p.id = pm.product_id 
		LEFT JOIN product_categories pc ON pc.product_id = p.id
		LEFT JOIN categories c ON c.id = pc.category_id
	GROUP BY 
		p.id, 
		p.product_code, 
		p.name, p.description, 
		pm.id, pm.media_url, pc.category_id, 
		c.name, vt1.variant_group_id, 
		vg1.name, 
		vt2.variant_group_id, 
		vg2.name
	HAVING p.product_code = $1;
	`

	err := r.db.SelectContext(ctx, &detail, qs, productCode)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (r *productRepository) FindProductVariants(ctx context.Context, productId int) ([]dto.GetProductVariant, error) {
	variant := make([]dto.GetProductVariant, 0)
	qs := `
	SELECT
		pv.price,
		pv.stock,
		pv.discount,
		pv.variant_type1_id,
		vt1.name AS variant_type1_name,
		pv.variant_type2_id,
		vt2.name AS variant_type2_name
	FROM
		product_variants pv 
		LEFT JOIN variant_types vt1 ON pv.variant_type1_id = vt1.id
		LEFT JOIN variant_types vt2 ON pv.variant_type2_id = vt2.id
	WHERE product_id = $1
	ORDER BY pv.id;
	`

	err := r.db.SelectContext(ctx, &variant, qs, productId)
	if err != nil {
		return nil, err
	}

	return variant, nil
}

func (r *productRepository) FirstProductByCode(ctx context.Context, code string) (*model.Product, error) {
	product := new(model.Product)
	err := r.db.GetContext(ctx, product, "SELECT * FROM products p WHERE p.product_code = $1", code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrProductNotFound
		}
		return nil, err
	}

	return product, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, payload dto.UpdateProductPayload, categories []interface{}, mediaType []string, accountId int64) error {
	count := 0
	qs1 := `
	UPDATE products
	SET name = $1, description = $2, weight = $3
	WHERE id = $4 AND seller_id = $5
	`

	qsa := `
	SELECT
		*
	FROM
		product_medias pm 
	WHERE
		product_id = $1
	ORDER BY id
	`

	qsb := `
	UPDATE product_medias
	SET media_url = $1
	WHERE id = (
		SELECT
			id
		FROM
			product_medias pm 
		WHERE
			product_id = $2
		ORDER BY id
		LIMIT 1
		OFFSET $3
	)
	`

	qsc := `
	DELETE FROM product_medias
	WHERE id IN (SELECT id FROM product_medias WHERE product_id = $1 ORDER BY id DESC LIMIT $2);
	`

	qsd := `
	INSERT INTO product_medias (media_url, media_type, product_id)
	VALUES
	($1,$2,$3)
	`

	qse := `
	UPDATE products
	SET thumbnail_url = $1
	WHERE id = $2
	`

	qs2 := `
	SELECT
		*
	FROM
		product_categories
	WHERE
		product_id = $1
	ORDER BY
		category_id
	`

	qs3 := `
	UPDATE product_categories
	SET category_id = $1
	WHERE id = (
		SELECT
			id
		FROM
			product_categories
		WHERE
			product_id = $2
		ORDER BY
			category_id
		LIMIT 1
		OFFSET $3
	)
	`

	qs4 := `
	DELETE FROM product_categories
	WHERE id IN (SELECT id FROM product_categories WHERE product_id = $1 ORDER BY id DESC LIMIT 1);
	`

	qs5 := `
	INSERT INTO product_categories (product_id, category_id)
	VALUES
	($1,$2)
	`

	qs6 := `
	UPDATE product_variants
	SET price = $1, stock = $2
	WHERE id = (
		SELECT
			pv.id
		FROM
			product_variants pv
			LEFT JOIN products p ON pv.product_id = p.id
			LEFT JOIN variant_types vt1 ON pv.variant_type1_id = vt1.id
			LEFT JOIN variant_types vt2 ON pv.variant_type2_id = vt2.id
		WHERE
			pv.product_id = $3 AND vt1.name = $4 AND vt2.name = $5
	)
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(qs1, payload.ProductName, payload.Description, payload.Weight, payload.ProductID, accountId)
	if err != nil {
		return err
	}

	rows1, err := tx.Queryx(qsa, payload.ProductID)
	if err != nil {
		return err
	}

	for rows1.Next() {
		count++
	}

	switch {
	case len(payload.ImageURL) < count:
		for idx, v := range payload.ImageURL {
			_, err = tx.Exec(qsb, v, payload.ProductID, idx)
			if err != nil {
				return err
			}
		}
		_, err = tx.Exec(qsc, payload.ProductID, count-len(payload.ImageURL))
		if err != nil {
			return err
		}
	case len(payload.ImageURL) > count:
		for idx, v := range payload.ImageURL {
			if idx >= count {
				_, err = tx.Exec(qsd, v, mediaType[idx], payload.ProductID)
				if err != nil {
					return err
				}
				continue
			}
			_, err = tx.Exec(qsb, v, payload.ProductID, idx)
			if err != nil {
				return err
			}
		}
	default:
		for idx, v := range payload.ImageURL {
			_, err = tx.Exec(qsb, v, payload.ProductID, idx)
			if err != nil {
				return err
			}
		}
	}

	_, err = tx.Exec(qse, payload.ImageURL[0], payload.ProductID)
	if err != nil {
		return err
	}

	count = 0
	rows2, err := tx.Queryx(qs2, payload.ProductID)
	if err != nil {
		return err
	}

	for rows2.Next() {
		count++
	}

	if categories[2] == 0 {
		categories = categories[:2]
	}

	switch {
	case len(categories) < count:
		for idx, v := range categories {
			_, err = tx.Exec(qs3, v, payload.ProductID, idx)
			if err != nil {
				return err
			}
		}
		_, err = tx.Exec(qs4, payload.ProductID)
		if err != nil {
			return err
		}
	case len(categories) > count:
		for idx, v := range categories {
			if idx == count {
				break
			}
			_, err = tx.Exec(qs3, v, payload.ProductID, idx)
			if err != nil {
				return err
			}
		}
		_, err = tx.Exec(qs5, payload.ProductID, categories[2])
		if err != nil {
			return err
		}
	default:
		for idx, v := range categories {
			_, err = tx.Exec(qs3, v, payload.ProductID, idx)
			if err != nil {
				return err
			}
		}
	}

	for idx, v := range payload.Variants {
		if payload.Variants[idx].VariantType1 == nil && payload.Variants[idx].VariantType2 == nil {
			_, err := tx.Exec(qs6, v.Price, v.Stock, payload.ProductID, constant.ProductVariantDefault, constant.ProductVariantDefault)
			if err != nil {
				return err
			}
		}

		if payload.Variants[idx].VariantType2 == nil {
			_, err := tx.Exec(qs6, v.Price, v.Stock, payload.ProductID, v.VariantType1, constant.ProductVariantDefault)
			if err != nil {
				return err
			}
		}

		_, err := tx.Exec(qs6, v.Price, v.Stock, payload.ProductID, v.VariantType1, v.VariantType2)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}
