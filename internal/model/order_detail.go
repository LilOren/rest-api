package model

import "database/sql"

type OrderDetail struct {
	ID            int64        `db:"id"`
	OrderID       int64        `db:"order_id"`
	ProductCode   string       `db:"product_code"`
	ProductName   string       `db:"product_name"`
	ThumbnailURL  string       `db:"thumbnail_url"`
	VariantName   string       `db:"variant_name"`
	SubTotalPrice float64      `db:"sub_total_price"`
	Quantity      int          `db:"quantity"`
	CreatedAt     sql.NullTime `db:"created_at"`
	UpdatedAt     sql.NullTime `db:"updated_at"`
	DeletedAt     sql.NullTime `db:"deleted_at"`
}
