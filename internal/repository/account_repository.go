package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/model"
	"github.com/lil-oren/rest/internal/shared"
)

type (
	AccountRepository interface {
		FirstByUsername(ctx context.Context, username string) (*model.Account, error)
		FirstByEmail(ctx context.Context, email string) (*model.Account, error)
		FirstById(ctx context.Context, id int64) (*model.Account, error)
		Create(ctx context.Context, e model.Account) error
		UpdateProfilePicture(ctx context.Context, accountID int64, photoURL string) error
		UpdatePassword(ctx context.Context, accountID int64, hashedPassword string) error
		UpdateWalletPin(ctx context.Context, accountID int64, pinHash string) error
	}
	accountRepository struct {
		db *sqlx.DB
		wr WalletRepository
	}
)

// UpdateWalletPin implements AccountRepository.
func (r *accountRepository) UpdateWalletPin(ctx context.Context, accountID int64, pinHash string) error {
	qs := `
	UPDATE accounts
	SET
		pin_hash = :pin_hash
	WHERE 
		id = :account_id
	`
	args := map[string]interface{}{
		"pin_hash":   pinHash,
		"account_id": accountID,
	}
	_, err := r.db.NamedExecContext(ctx, qs, args)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProfilePicture implements AccountRepository.
func (r *accountRepository) UpdateProfilePicture(ctx context.Context, accountID int64, photoURL string) error {

	qs := `
	UPDATE accounts a
	SET
		profile_picture_url = :photo_url
	WHERE
		a.id = :account_id
	`
	args := map[string]interface{}{
		"photo_url":  photoURL,
		"account_id": accountID,
	}

	res, err := r.db.NamedExecContext(ctx, qs, args)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return shared.ErrAccountNotFound
	}

	return nil
}

// FirstById implements AccountRepository.
func (r *accountRepository) FirstById(ctx context.Context, id int64) (*model.Account, error) {
	e := new(model.Account)

	rows := r.db.QueryRowxContext(ctx, "SELECT * FROM accounts WHERE id = $1", id)
	if err := rows.Err(); err != nil {
		return nil, err
	}

	err := rows.StructScan(e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// FirstByEmail implements AccountRepository.
func (r *accountRepository) FirstByEmail(ctx context.Context, email string) (*model.Account, error) {
	e := new(model.Account)

	rows := r.db.QueryRowxContext(ctx, "SELECT * FROM accounts WHERE email = $1", email)
	if err := rows.Err(); err != nil {
		return nil, err
	}

	err := rows.StructScan(e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// Create implements AccountRepository.
func (r *accountRepository) Create(ctx context.Context, e model.Account) error {
	qs := `
	INSERT INTO accounts (
		username, 
		email, 
		password_hash
		) VALUES (
			$1,
			$2,
			$3
		)
	RETURNING id
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	userID := new(int64)
	if err = tx.Get(userID, qs, e.Username, e.Email, e.PasswordHash); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := r.wr.TxCreateWalletAfterRegister(tx, *userID); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// FirstByUsername implements AccountRepository.
func (r *accountRepository) FirstByUsername(ctx context.Context, username string) (*model.Account, error) {
	e := new(model.Account)

	row := r.db.QueryRowContext(ctx, "SELECT * FROM accounts WHERE username = $1", username)
	if err := row.Scan(e); err != nil {
		return nil, err
	}

	return e, nil
}

func (r *accountRepository) UpdatePassword(ctx context.Context, accountID int64, hashedPassword string) error {
	query := `
		UPDATE accounts a
		SET	password_hash = :hashed_password
		WHERE a.id = :account_id
	`
	args := map[string]interface{}{
		"hashed_password": hashedPassword,
		"account_id":      accountID,
	}

	res, err := r.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return shared.ErrAccountNotFound
	}

	return nil
}

func NewAccountRepository(db *sqlx.DB, wr WalletRepository) AccountRepository {
	return &accountRepository{
		db: db,
		wr: wr,
	}
}
