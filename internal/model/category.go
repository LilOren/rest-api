package model

import "database/sql"

type Category struct {
	ID             int64         `db:"id"`
	Name           string        `db:"name"`
	Level          uint8         `db:"level"`
	ImageUrl       string        `db:"image_url"`
	ParentCategory sql.NullInt64 `db:"parent_category"`
	CreatedAt      sql.NullTime  `db:"created_at"`
	UpdatedAt      sql.NullTime  `db:"updated_at"`
	DeletedAt      sql.NullTime  `db:"deleted_at"`
}
