package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
)

type (
	AccountAddressRepository interface {
		CreateAddress(ctx context.Context, payload model.AccountAddresses, length int, accountId int) error
		FirstShopAddressById(ctx context.Context, accountId int) (*model.AccountAddresses, error)
		FindAddressById(ctx context.Context, accountId int) ([]model.AccountAddresses, error)
		UpdateDefaultAddress(ctx context.Context, accountId int, id int) error
		FindDetailsAddressById(ctx context.Context, accountId int) ([]dto.AccountDetailsAddress, error)
		FirstShopAddressByShopID(ctx context.Context, shopID int64) (*model.AccountAddresses, error)
		FirstByID(ctx context.Context, id int64) (*model.AccountAddresses, error)
		UpdateAddressByID(ctx context.Context, payload model.AccountAddresses) error
	}
	accountAddressRepository struct {
		db *sqlx.DB
	}
)

// UpdateAddressByID implements AccountAddressRepository.
func (r *accountAddressRepository) UpdateAddressByID(ctx context.Context, payload model.AccountAddresses) error {

	qs := `
	UPDATE account_addresses
	SET
		receiver_name = :receiver_name,
		receiver_phone_number = :receiver_phone_number,
		detail = :detail,
		is_shop = :is_shop,
		is_default = :is_default,
		province_id = :province_id,
		district_id = :district_id,
		postal_code = :postal_code,
		account_id = :account_id,
		updated_at = NOW()
	WHERE
		id = :id
	`

	_, err := r.db.NamedExecContext(ctx, qs, payload)
	if err != nil {
		return err
	}

	return nil
}

// FirstByID implements AccountAddressRepository.
func (r *accountAddressRepository) FirstByID(ctx context.Context, id int64) (*model.AccountAddresses, error) {
	qs := `
	SELECT
	aa.*
	FROM account_addresses aa
	WHERE aa.id = $1
	`
	row := r.db.QueryRowxContext(ctx, qs, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	m := new(model.AccountAddresses)
	if err := row.StructScan(m); err != nil {
		return nil, err
	}

	return m, nil
}

// FirstShopAddressByShopID implements AccountAddressRepository.
func (r *accountAddressRepository) FirstShopAddressByShopID(ctx context.Context, shopID int64) (*model.AccountAddresses, error) {
	qs := `
	SELECT 
	aa.*
	FROM shops s
	LEFT JOIN accounts a
		ON a.id = s.account_id
	LEFT JOIN account_addresses aa
		ON a.id = aa.account_id
	WHERE 
		aa.is_shop = TRUE AND
		s.id = $1
	`

	m := new(model.AccountAddresses)
	row := r.db.QueryRowxContext(ctx, qs, shopID)
	if err := row.Err(); err != nil {
		return nil, err
	}

	if err := row.StructScan(m); err != nil {
		return nil, err
	}

	return m, nil
}

// Create implements AccountAddressRepository.
func (aar *accountAddressRepository) CreateAddress(ctx context.Context, address model.AccountAddresses, length int, accountId int) error {
	qs := `
	INSERT INTO account_addresses (
		receiver_name,
		receiver_phone_number,
		detail, 
		is_shop,
		is_default,
		province_id,
		district_id,
		postal_code,
		account_id,
		updated_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10
		)
	`

	if length == 0 {
		_, err := aar.db.ExecContext(ctx, qs, address.ReceiverName, address.ReceiverPhoneNumber, address.Detail, false, true, address.ProvinceId, address.DistrictId, address.PostalCode, accountId, time.Now())
		if err != nil {
			return err
		}
	}

	if length != 0 {
		_, err := aar.db.ExecContext(ctx, qs, address.ReceiverName, address.ReceiverPhoneNumber, address.Detail, false, false, address.ProvinceId, address.DistrictId, address.PostalCode, accountId, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

func (aar *accountAddressRepository) FirstShopAddressById(ctx context.Context, accountId int) (*model.AccountAddresses, error) {
	address := new(model.AccountAddresses)
	qs := `
	SELECT 
		* 
	FROM
		account_addresses
	WHERE
		account_id = $1
	ORDER BY
		is_shop DESC
	`

	row := aar.db.QueryRowxContext(ctx, qs, accountId)
	if err := row.Err(); err != nil {
		return nil, err
	}

	err := row.StructScan(address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (aar *accountAddressRepository) FindAddressById(ctx context.Context, accountId int) ([]model.AccountAddresses, error) {
	addressList := make([]model.AccountAddresses, 0)
	qs := `
	SELECT 
		* 
	FROM
		account_addresses
	WHERE
		account_id = $1
	`

	rows, err := aar.db.QueryxContext(ctx, qs, accountId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := new(model.AccountAddresses)
		rows.StructScan(temp)
		addressList = append(addressList, *temp)
	}

	return addressList, nil
}

func (aar *accountAddressRepository) FindDetailsAddressById(ctx context.Context, accountId int) ([]dto.AccountDetailsAddress, error) {
	addressList := make([]dto.AccountDetailsAddress, 0)
	qs := `
	SELECT
		aa.id,
		aa.receiver_name,
		aa.detail,
		aa.postal_code,
		aa.receiver_phone_number
	FROM
		account_addresses aa
	WHERE
		account_id = $1
	ORDER BY
		aa.is_default DESC
	`

	rows, err := aar.db.QueryxContext(ctx, qs, accountId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := new(dto.AccountDetailsAddress)
		if err := rows.StructScan(temp); err != nil {
			return nil, err
		}
		addressList = append(addressList, *temp)
	}

	return addressList, nil
}

func (aar *accountAddressRepository) UpdateDefaultAddress(ctx context.Context, accountId int, id int) error {
	qs1 := `
	UPDATE account_addresses
	SET is_default = $1, updated_at = $2
	WHERE account_id = $3 AND id = $4
	`

	qs2 := `
	UPDATE account_addresses
	SET is_default = $1, updated_at = $2
	WHERE account_id = $3 AND id != $4
	`

	tx, err := aar.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(qs1, true, time.Now(), accountId, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(qs2, false, time.Now(), accountId, id)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func NewAccountAddressRepository(db *sqlx.DB) AccountAddressRepository {
	return &accountAddressRepository{
		db: db,
	}
}
