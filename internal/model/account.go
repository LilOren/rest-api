package model

import "database/sql"

type Account struct {
	ID                int64          `db:"id"`
	Username          string         `db:"username"`
	Email             string         `db:"email"`
	PasswordHash      sql.NullString `db:"password_hash"`
	Fullname          sql.NullString `db:"full_name"`
	PhoneNumber       sql.NullString `db:"phone_number"`
	Gender            sql.NullString `db:"gender"`
	BirthDate         sql.NullTime   `db:"birth_date"`
	IsSeller          bool           `db:"is_seller"`
	ProfilePictureURL sql.NullString `db:"profile_picture_url"`
	PinHash           sql.NullString `db:"pin_hash"`
	CreatedAt         sql.NullTime   `db:"created_at"`
	UpdatedAt         sql.NullTime   `db:"updated_at"`
	DeletedAt         sql.NullTime   `db:"deleted_at"`
}
