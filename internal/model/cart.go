package model

import "database/sql"

type Cart struct {
	ID               int64        `db:"id"`
	ProductVariantId int64        `db:"product_variant_id"`
	AccountId        int64        `db:"account_id"`
	SellerId         int64        `db:"seller_id"`
	Quantity         int          `db:"quantity"`
	IsChecked        bool         `db:"is_checked"`
	CreatedAt        sql.NullTime `db:"created_at"`
	UpdatedAt        sql.NullTime `db:"updated_at"`
	DeletedAt        sql.NullTime `db:"deleted_at"`
}
