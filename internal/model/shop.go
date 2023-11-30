package model

import "database/sql"

type Shop struct {
	ID        int64          `db:"id"`
	Name      sql.NullString `db:"name"`
	AccountId int            `db:"account_id"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
}
