package model

import (
	"database/sql"

	"github.com/lil-oren/rest/internal/constant"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID           int64                     `db:"id"`
	Amount       decimal.Decimal           `db:"amount"`
	Title        constant.TransactionTitle `db:"title"`
	FromWalletID sql.NullInt64             `db:"from_wallet_id"`
	ToWalletID   int64                     `db:"to_wallet_id"`
	CreatedAt    sql.NullTime              `db:"created_at"`
}
