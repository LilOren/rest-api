package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/dto"
)

type (
	SellerPageRepository interface {
		FirstSellerDetail(ctx context.Context, shopName string) (*dto.GetSellerDetail, error)
		FindSellerProductBySearchTerm(ctx context.Context, payload dto.SearchSellerProductPayload, sellerId int64) ([]dto.SearchSellerProductResponseDb, error)
		CountProductBySearchTerm(ctx context.Context, payload dto.SearchSellerProductPayload, sellerId int64) (*int, error)
		FindCategoryBySellerId(ctx context.Context, sellerId int64) ([]string, error)
	}
	sellerPageRepository struct {
		db *sqlx.DB
	}
)

func (spr *sellerPageRepository) FirstSellerDetail(ctx context.Context, shopName string) (*dto.GetSellerDetail, error) {
	shopDetail := new(dto.GetSellerDetail)
	qs := `
	SELECT
		s.account_id AS seller_id,
		s.name AS shop_name,
		count(*) AS product_counts,
		DATE_PART('year', AGE(now(), s.created_at)) AS years
	FROM
		product_variants pv
		LEFT JOIN products p ON pv.product_id = p.id
		LEFT JOIN shops s ON p.seller_id = s.id
	GROUP BY 
		s.account_id,
		s.name, 
		s.created_at
	HAVING 
		s.name = $1;
	`

	err := spr.db.GetContext(ctx, shopDetail, qs, shopName)
	if err != nil {
		return nil, err
	}
	return shopDetail, nil
}

func (spr *sellerPageRepository) FindSellerProductBySearchTerm(ctx context.Context, payload dto.SearchSellerProductPayload, sellerId int64) ([]dto.SearchSellerProductResponseDb, error) {
	qs := `
	SELECT
		p.product_code AS product_code,
		p.thumbnail_url as thumbnail_url,
		p.name AS product_name,
		pv.price as base_price,
		pv.discount as discount,
		(pv.price - (pv.price * pv.discount / 100)) AS discount_price,
		d.name AS district_name,
		(CASE 
			WHEN mp.count_purchased IS NULL THEN 0
			ELSE mp.count_purchased
		END) AS count_purchased
	FROM
		(
			SELECT 
				DISTINCT ON (pc.product_id)
				p.*,
				pc.category_id,
				c."name" AS category_name
			FROM
				products p
				LEFT JOIN product_categories pc ON pc.product_id = p.id
				LEFT JOIN categories c ON pc.category_id = c.id 
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
	%s
	ORDER BY
		%s
	OFFSET :offset
	LIMIT 20
	`

	start := 0
	searchTerm := "%"

	if payload.SearchTerm != "" {
		searchTerm = "%" + payload.SearchTerm + "%"
	}

	if payload.Page > 1 {
		start = (payload.Page - 1) * 20
	}

	args := map[string]interface{}{
		"offset":      start,
		"search_term": searchTerm,
	}

	condition := "WHERE p.seller_id = %d"
	condition = fmt.Sprintf(condition, sellerId)
	qs = fmt.Sprintf(qs, "%s", condition, "%s")

	if payload.CategoryName != "" {
		condition := "AND c.name = '" + "%s" + "'"
		condition = fmt.Sprintf(condition, payload.CategoryName)
		qs = fmt.Sprintf(qs, condition, "%s")
	} else {
		qs = fmt.Sprintf(qs, "", "%s")
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
				qs = fmt.Sprintf(qs, "count_purchased DESC")
				break
			}

			qs = fmt.Sprintf(qs, "count_purchased ASC")
		}
	}

	stmt, err := spr.db.PrepareNamedContext(ctx, qs)
	if err != nil {
		return nil, err
	}

	m := make([]dto.SearchSellerProductResponseDb, 0)
	if err := stmt.SelectContext(ctx, &m, args); err != nil {
		return nil, err
	}

	return m, nil
}

func (spr *sellerPageRepository) CountProductBySearchTerm(ctx context.Context, payload dto.SearchSellerProductPayload, sellerId int64) (*int, error) {
	qs := `
	SELECT
		p.product_code AS product_code,
		p.thumbnail_url as thumbnail_url,
		p.name AS product_name,
		pv.price as base_price,
		pv.discount as discount,
		(pv.price - (pv.price * pv.discount / 100)) AS discount_price,
		d.name AS district_name,
		(CASE 
			WHEN mp.count_purchased IS NULL THEN 0
			ELSE mp.count_purchased
		END) AS count_purchased
	FROM
		(
			SELECT 
				DISTINCT ON (pc.product_id)
				p.*,
				pc.category_id,
				c."name" AS category_name
			FROM
				products p
				LEFT JOIN product_categories pc ON pc.product_id = p.id
				LEFT JOIN categories c ON pc.category_id = c.id 
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
	%s
	ORDER BY
		%s
	`

	start := 0
	searchTerm := "%"

	if payload.SearchTerm != "" {
		searchTerm = "%" + payload.SearchTerm + "%"
	}

	if payload.Page > 1 {
		start = (payload.Page - 1) * 20
	}

	args := map[string]interface{}{
		"offset":      start,
		"search_term": searchTerm,
	}

	condition := "WHERE p.seller_id = %d"
	condition = fmt.Sprintf(condition, sellerId)
	qs = fmt.Sprintf(qs, "%s", condition, "%s")

	if payload.CategoryName != "" {
		condition := "AND c.name = '" + "%s" + "'"
		condition = fmt.Sprintf(condition, payload.CategoryName)
		qs = fmt.Sprintf(qs, condition, "%s")
	} else {
		qs = fmt.Sprintf(qs, "", "%s")
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
				qs = fmt.Sprintf(qs, "count_purchased DESC")
				break
			}

			qs = fmt.Sprintf(qs, "count_purchased ASC")
		}
	}

	stmt, err := spr.db.PrepareNamedContext(ctx, qs)
	if err != nil {
		return nil, err
	}

	m := make([]dto.SearchSellerProductResponseDb, 0)
	if err := stmt.SelectContext(ctx, &m, args); err != nil {
		return nil, err
	}

	res := len(m)

	return &res, nil
}

func (spr *sellerPageRepository) FindCategoryBySellerId(ctx context.Context, sellerId int64) ([]string, error) {
	categoryList := make([]string, 0)
	qs := `
	SELECT
		t.name
	FROM (
		SELECT
			pc.product_id,
			pc.category_id,
			c.name,
			row_number() OVER (PARTITION BY pc.product_id ORDER BY pc.category_id DESC) AS rn,
			count(pc.product_id) OVER (PARTITION BY pc.product_id) cn 
		FROM
			product_categories pc 
			LEFT JOIN products p ON p.id = pc.product_id
			LEFT JOIN categories c ON c.id = pc.category_id 
		WHERE p.seller_id = $1
		ORDER BY pc.id
	) AS t
	WHERE cn > 1 AND rn = 1
	GROUP BY t.name
	`

	rows, err := spr.db.QueryxContext(ctx, qs, sellerId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := new(string)
		if err := rows.Scan(temp); err != nil {
			return nil, err
		}
		categoryList = append(categoryList, *temp)
	}

	return categoryList, nil
}

func NewSellerPageRepository(db *sqlx.DB) SellerPageRepository {
	return &sellerPageRepository{
		db: db,
	}
}
