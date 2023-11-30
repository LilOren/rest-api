package model

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID        int64           `db:"id"`
	Balance   decimal.Decimal `db:"balance"`
	IsActive  bool            `db:"is_active"`
	Category  string          `db:"category"`
	AccountId int64           `db:"account_id"`
	CreatedAt sql.NullTime    `db:"created_at"`
	UpdatedAt sql.NullTime    `db:"updated_at"`
	DeletedAt sql.NullTime    `db:"deleted_at"`
}
