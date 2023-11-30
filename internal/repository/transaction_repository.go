package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
)

type (
	TransactionRepository interface {
		FirstTransactionByID(ctx context.Context, id int64) (*model.Transaction, error)
		CreateTransaction(tx *sqlx.Tx, transaction *model.Transaction) (*int64, error)
		CreateTopupTransaction(tx *sqlx.Tx, transaction *model.Transaction) (*int64, error)
		FindTransactionByAccountID(ctx context.Context, walletID int64, startDate, endDate time.Time, offset int) ([]dto.ListWalletHistoryItemFromDB, error)
		CountTransactionByAccountID(ctx context.Context, walletID int64, startDate, endDate time.Time) (*int64, error)
	}
	transactionRepository struct {
		db *sqlx.DB
	}
)

// CountTransactionByAccountID implements TransactionRepository.
func (r *transactionRepository) CountTransactionByAccountID(ctx context.Context, walletID int64, startDate time.Time, endDate time.Time) (*int64, error) {
	qs := `
	SELECT
		COUNT(1) AS total_count
	FROM
		transactions t
	FULL OUTER JOIN orders o ON
		o.transaction_id = t.id
	LEFT JOIN shops s ON
		s.account_id = o.seller_id
	LEFT JOIN wallets w ON
		w.id = t.to_wallet_id
		OR w.id = t.from_wallet_id
	WHERE
			w.id = :wallet_id
		AND
			t.created_at >= :start_date
		AND
		t.created_at <= :end_date`

	args := map[string]interface{}{
		"wallet_id":  walletID,
		"start_date": startDate.Format(constant.DateLayoutISO),
		"end_date":   endDate.Format(constant.DateLayoutISO),
	}

	rows, err := r.db.NamedQueryContext(ctx, qs, args)
	if err != nil {
		return nil, err
	}

	rows.Next()
	count := new(int64)
	if err := rows.Scan(count); err != nil {
		return nil, err
	}

	return count, nil

}

// FindTransactionByAccountID implements TransactionRepository.
func (r *transactionRepository) FindTransactionByAccountID(
	ctx context.Context,
	walletID int64,
	startDate, endDate time.Time,
	offset int,
) ([]dto.ListWalletHistoryItemFromDB, error) {
	res := make([]dto.ListWalletHistoryItemFromDB, 0)

	qs := `
	SELECT
		t.title,
		t.amount,
		t.to_wallet_id = w.id AS is_debit,
		o.id AS order_id,
		s.name AS shop_name,
		t.created_at as date
	FROM
		transactions t
	FULL OUTER JOIN orders o ON
		o.transaction_id = t.id
	LEFT JOIN shops s ON
		s.account_id = o.seller_id
	LEFT JOIN wallets w ON
		w.id = t.to_wallet_id
		OR w.id = t.from_wallet_id
	WHERE
		w.id = :wallet_id AND
		t.created_at >= :start_date AND
		t.created_at <= :end_date
	ORDER BY
		t.created_at DESC
	OFFSET :offset
	LIMIT 10
	`

	args := map[string]interface{}{
		"wallet_id":  walletID,
		"start_date": startDate.Format(constant.DateLayoutISO),
		"end_date":   endDate.Format(constant.DateLayoutISO),
		"offset":     offset,
	}

	rows, err := r.db.NamedQueryContext(ctx, qs, args)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := new(dto.ListWalletHistoryItemFromDB)
		if err := rows.StructScan(temp); err != nil {
			return nil, err
		}

		res = append(res, *temp)
	}

	return res, nil
}

// CreateTransaction implements TransactionRepository.
func (*transactionRepository) CreateTransaction(tx *sqlx.Tx, transaction *model.Transaction) (*int64, error) {
	var id int64
	query := `
	INSERT INTO transactions 
	(
		amount, 
		title, 
		to_wallet_id, 
		from_wallet_id
	)
	VALUES
	($1, $2, $3, $4)
	RETURNING(id)`
	err := tx.QueryRowx(query, transaction.Amount, transaction.Title, transaction.ToWalletID, transaction.FromWalletID.Int64).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// CreateTopupTransaction implements TransactionRepository.
func (*transactionRepository) CreateTopupTransaction(tx *sqlx.Tx, transaction *model.Transaction) (*int64, error) {
	id := new(int64)
	qs := `
	INSERT INTO transactions
	(
		amount,
		title,
		to_wallet_id
	) VALUES
	(
		:amount, 
		:title, 
		:to_wallet_id
	)
	RETURNING (id)
	`

	rows, err := tx.NamedQuery(qs, transaction)
	if err != nil {
		return nil, err
	}

	rows.Next()
	if err := rows.Scan(id); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := rows.Close(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return id, nil
}

func (r *transactionRepository) FirstTransactionByID(ctx context.Context, id int64) (*model.Transaction, error) {
	query := `SELECT * FROM transactions t WHERE t.id = $1`
	transaction := new(model.Transaction)
	err := r.db.GetContext(ctx, transaction, query, id)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}
