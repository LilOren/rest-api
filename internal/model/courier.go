package model

import "database/sql"

type Courier struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	Code        string `db:"code"`
	ServiceName string `db:"service_name"`
	Description string `db:"description"`
	ImageUrl    string `db:"image_url"`

	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
