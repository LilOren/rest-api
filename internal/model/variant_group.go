package model

import "database/sql"

type VariantGroup struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	ProductID int64        `db:"product_id"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
