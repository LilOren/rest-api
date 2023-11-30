package model

import "database/sql"

type ChangedEmail struct {
	ID        int64        `db:"id"`
	AccountID int64        `db:"account_id"`
	Email     string       `db:"email"`
	CreatedAt sql.NullTime `db:"created_at"`
}
