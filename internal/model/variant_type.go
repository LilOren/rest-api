package model

import "database/sql"

type VariantType struct {
	ID             int64        `db:"id"`
	Name           string       `db:"name"`
	VariantGroupID int64        `db:"variant_group_id"`
	CreatedAt      sql.NullTime `db:"created_at"`
	UpdatedAt      sql.NullTime `db:"updated_at"`
	DeletedAt      sql.NullTime `db:"deleted_at"`
}
