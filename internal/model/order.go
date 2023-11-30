package model

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID              int64           `db:"id"`
	Status          string          `db:"status"`
	EAT             sql.NullTime    `db:"estimated_time_arrival"`
	DeliveryCost    decimal.Decimal `db:"delivery_cost"`
	CourierId       int64           `db:"courier_id"`
	SellerId        int64           `db:"seller_id"`
	BuyerId         int64           `db:"buyer_id"`
	TransactionId   int64           `db:"transaction_id"`
	PromotionName   sql.NullString  `db:"promotion_name"`
	PromotionAmount sql.NullFloat64 `db:"promotion_amount"`
	CreatedAt       sql.NullTime    `db:"created_at"`
	UpdatedAt       sql.NullTime    `db:"updated_at"`
	DeletedAt       sql.NullTime    `db:"deleted_at"`
}
