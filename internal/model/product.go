package model

import "database/sql"

type Product struct {
	ID           int64        `db:"id"`
	Name         string       `db:"name"`
	ProductCode  string       `db:"product_code"`
	Description  string       `db:"description"`
	ThumbnailUrl string       `db:"thumbnail_url"`
	SellerID     int64        `db:"seller_id"`
	Weight       int          `db:"weight"`
	CreatedAt    sql.NullTime `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}
