package model

import "database/sql"

type AccountAddresses struct {
	ID                  int64        `db:"id"`
	ReceiverName        string       `db:"receiver_name"`
	ReceiverPhoneNumber string       `db:"receiver_phone_number"`
	Detail              string       `db:"detail"`
	IsShop              bool         `db:"is_shop"`
	IsDefault           bool         `db:"is_default"`
	ProvinceId          int64        `db:"province_id"`
	DistrictId          int64        `db:"district_id"`
	PostalCode          string       `db:"postal_code"`
	AccountId           int64        `db:"account_id"`
	CreatedAt           sql.NullTime `db:"created_at"`
	UpdatedAt           sql.NullTime `db:"updated_at"`
	DeletedAt           sql.NullTime `db:"deleted_at"`
}
