package model

import "database/sql"

type ProductMedia struct {
	ID        int64        `db:"id"`
	MediaUrl  string       `db:"media_url"`
	MediaType string       `db:"media_type"`
	ProductID int64        `db:"product_id"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
