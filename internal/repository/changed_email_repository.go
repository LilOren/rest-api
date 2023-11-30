package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/model"
)

type (
	ChangedEmailRepository interface {
		FirstByEmail(ctx context.Context, email string) (*model.ChangedEmail, error)
		Create(ctx context.Context, userID int64, newEmail, oldEmail string) error
	}
	changedEmailRepository struct {
		db *sqlx.DB
	}
)

// Create implements ChangedEmailRepository.
func (r *changedEmailRepository) Create(ctx context.Context, userID int64, newEmail, oldEmail string) error {
	qs := `
	INSERT INTO changed_emails
	(
		account_id,
		email
	)
	VALUES
	($1, $2)
	`

	qs2 := `
	UPDATE accounts
	SET
		email = $1
	WHERE id = $2
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(qs, userID, oldEmail)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	_, err = tx.Exec(qs2, newEmail, userID)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// FirstByID implements ChangedEmailRepository.
func (r *changedEmailRepository) FirstByEmail(ctx context.Context, email string) (*model.ChangedEmail, error) {
	qs := `
	SELECT 
	* 
	FROM changed_emails
	WHERE
		email = $1
	LIMIT 1
	`

	row := r.db.QueryRowxContext(ctx, qs, email)
	if err := row.Err(); err != nil {
		return nil, err
	}

	m := new(model.ChangedEmail)
	if err := row.StructScan(m); err != nil {
		return nil, err
	}

	return m, nil
}

func NewChangedEmailRepository(db *sqlx.DB) ChangedEmailRepository {
	return &changedEmailRepository{
		db: db,
	}
}
