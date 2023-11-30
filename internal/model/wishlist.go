package model

import "database/sql"

type Wishlist struct {
	ID        int64        `db:"id"`
	AccountID int64        `db:"account_id"`
	ProductID int64        `db:"product_id"`
	CreatedAt sql.NullTime `db:"created_at"`
}
