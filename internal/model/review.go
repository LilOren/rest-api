package model

import "database/sql"

type Review struct {
	ID          int64        `db:"id"`
	Rating      int          `db:"rating"`
	Comment     string       `db:"comment"`
	AccountID   int64        `db:"account_id"`
	ProductCode string       `db:"product_code"`
	CreatedAt   sql.NullTime `db:"created_at"`
}
