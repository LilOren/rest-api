package model

type District struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	ProvinceID int64  `db:"province_id"`
	PostalCode string `db:"postal_code"`
}
