package model

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type ProductVariant struct {
	ID             int64           `db:"id"`
	Price          decimal.Decimal `db:"price"`
	Stock          uint32          `db:"stock"`
	Discount       float32         `db:"discount"`
	ProductID      int64           `db:"product_id"`
	VariantType1ID int64           `db:"variant_type1_id"`
	VariantType2ID int64           `db:"variant_type2_id"`
	CreatedAt      sql.NullTime    `db:"created_at"`
	UpdatedAt      sql.NullTime    `db:"updated_at"`
	DeletedAt      sql.NullTime    `db:"deleted_at"`
}
