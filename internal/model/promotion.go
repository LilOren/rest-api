package model

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Promotion struct {
	ID           int64           `db:"id"`
	Name         string          `db:"name"`
	ExactPrice   sql.NullFloat64 `db:"exact_price"`
	Percentage   sql.NullFloat64 `db:"percentage"`
	MinimumSpend decimal.Decimal `db:"minimum_spend"`
	Quota        int             `db:"quota"`
	ShopID       int64           `db:"shop_id"`
	StartedAt    time.Time       `db:"started_at"`
	ExpiredAt    time.Time       `db:"expired_at"`
	CreatedAt    time.Time       `db:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at"`
	DeletedAt    sql.NullTime    `db:"deleted_at"`
}
