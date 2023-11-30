package model

import "database/sql"

type ProductCategory struct {
	ID         int64        `db:"id"`
	ProductID  int64        `db:"product_id"`
	CategoryID int64        `db:"category_id"`
	CreatedAt  sql.NullTime `db:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
	DeletedAt  sql.NullTime `db:"deleted_at"`
}
