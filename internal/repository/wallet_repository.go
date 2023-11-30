package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/shopspring/decimal"
)

type (
	WalletRepository interface {
		CreateWallet(ctx context.Context, accountId int64) error
		TxCreateWalletAfterRegister(tx *sqlx.Tx, accountID int64) error
		ActivateShopWallet(ctx context.Context, accountId int64) error
		FirstActiveWalletByAccountID(ctx context.Context, accountId int64, category constant.WalletType) (*model.Wallet, error)
		ActivatePersonalAndTemporaryWallet(ctx context.Context, accountId int64, pinHash string) error
		WithdrawShopUser(ctx context.Context, accountId int64, transaction *model.Transaction) error
		Topup(ctx context.Context, accountId int64, transaction *model.Transaction) error
		FirstShopWalletBalanceBySellerID(ctx context.Context, sellerID int64) (*dto.ShopWalletBalanceResponse, error)
	}
	walletRepository struct {
		db *sqlx.DB
		tr TransactionRepository
	}
)

// TxCreateWalletAfterRegister implements WalletRepository.
func (r *walletRepository) TxCreateWalletAfterRegister(tx *sqlx.Tx, accountID int64) error {
	qs1 := `
	INSERT INTO wallets (balance, is_active, category, account_id) VALUES
	($1, $2, $3, $4)
	`
	qs2 := `
	INSERT INTO wallets (balance, is_active, category, account_id) VALUES
	($1, $2, $3, $4)
	`
	qs3 := `
	INSERT INTO wallets (balance, is_active, category, account_id) VALUES
	($1, $2, $3, $4)
	`

	_, err := tx.Exec(qs1, 0, false, constant.UserWalletType, accountID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(qs2, 0, false, constant.TempWalletType, accountID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(qs3, 0, false, constant.ShopWalletType, accountID)
	if err != nil {
		return err
	}

	return nil
}

func (r *walletRepository) CreateWallet(ctx context.Context, accountId int64) error {
	qs1 := `
	INSERT INTO wallets (balance, is_active, category, account_id) VALUES
	($1, $2, $3, $4)
	`
	qs2 := `
	INSERT INTO wallets (balance, is_active, category, account_id) VALUES
	($1, $2, $3, $4)
	`
	qs3 := `
	INSERT INTO wallets (balance, is_active, category, account_id) VALUES
	($1, $2, $3, $4)
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(qs1, 0, false, constant.UserWalletType, accountId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(qs2, 0, false, constant.TempWalletType, accountId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(qs3, 0, false, constant.ShopWalletType, accountId)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *walletRepository) ActivateShopWallet(ctx context.Context, accountId int64) error {
	qs := `
	UPDATE wallets
	SET 
		is_active = $1
	WHERE 
		category = $2 AND 
		account_id = $3
	`

	_, err := r.db.ExecContext(ctx, qs, true, constant.ShopWalletType, accountId)
	if err != nil {
		return err
	}

	return nil
}

func (r *walletRepository) ActivatePersonalAndTemporaryWallet(ctx context.Context, accountId int64, pinHash string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	qs1 := `
	UPDATE accounts
	SET
		pin_hash = $1
	WHERE 
		id = $2
	`

	_, err = tx.Exec(qs1, pinHash, accountId)
	if err != nil {
		tx.Rollback()
		return err
	}

	qs2 := `
	UPDATE wallets
	SET 
		is_active = TRUE	
	WHERE 
		account_id = $1 AND
		(category = $2 OR category = $3)
	`

	_, err = tx.Exec(qs2, accountId, constant.UserWalletType, constant.TempWalletType)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *walletRepository) FirstActiveWalletByAccountID(ctx context.Context, accountId int64, category constant.WalletType) (*model.Wallet, error) {
	eWallet := new(model.Wallet)
	query := `
	SELECT 
		*
	FROM
		wallets
	WHERE
		account_id = $1 AND
		category = $2 AND
		is_active
	`
	err := r.db.GetContext(ctx, eWallet, query, accountId, category)
	if err != nil {
		return nil, err
	}
	return eWallet, nil
}

func (r *walletRepository) WithdrawShopUser(ctx context.Context, accountId int64, transaction *model.Transaction) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = r.tr.CreateTransaction(tx, transaction)
	if err != nil {
		return err
	}

	wallet1 := new(model.Wallet)
	err = tx.Get(wallet1,
		"SELECT * FROM wallets w WHERE w.account_id = $1 AND w.category = $2 AND w.is_active FOR UPDATE",
		accountId, constant.ShopWalletType)
	if err != nil {
		return err
	}

	wallet2 := new(model.Wallet)
	err = tx.Get(wallet2,
		"SELECT * FROM wallets w WHERE w.account_id = $1 AND w.category = $2 AND w.is_active FOR UPDATE",
		accountId, constant.UserWalletType)
	if err != nil {
		return err
	}

	query1 := `
	UPDATE wallets
	SET balance = balance-$1
	WHERE account_id = $2 AND category = $3 AND is_active
	`

	query2 := `
	UPDATE wallets
	SET balance = balance+$1
	WHERE account_id = $2 AND category = $3 AND is_active
	`

	row1, err := tx.Exec(query1, transaction.Amount, accountId, constant.ShopWalletType)
	if err != nil {
		return err
	}
	aff1, _ := row1.RowsAffected()
	if aff1 == 0 {
		return shared.ErrUpdateInactiveWallet
	}

	row2, err := tx.Exec(query2, transaction.Amount, accountId, constant.UserWalletType)
	if err != nil {
		return err
	}
	aff2, _ := row2.RowsAffected()
	if aff2 == 0 {
		return shared.ErrUpdateInactiveWallet
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *walletRepository) Topup(ctx context.Context, accountId int64, transaction *model.Transaction) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = r.tr.CreateTopupTransaction(tx, transaction)
	if err != nil {
		return err
	}

	wallet := new(model.Wallet)
	err = tx.Get(wallet,
		"SELECT * FROM wallets w WHERE w.account_id = $1 AND w.category = $2 AND w.is_active FOR UPDATE",
		accountId, constant.UserWalletType)
	if err != nil {
		return err
	}

	query := `
	UPDATE wallets
	SET balance = balance+$1
	WHERE account_id = $2 AND category = $3 AND is_active
	`

	row, err := tx.Exec(query, transaction.Amount.InexactFloat64(), accountId, constant.UserWalletType)
	if err != nil {
		return err
	}
	aff, _ := row.RowsAffected()
	if aff == 0 {
		return shared.ErrUpdateInactiveWallet
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func transferAndRefund(tx *sqlx.Tx, accountId int64, totalPrice decimal.Decimal, wallet1Type, wallet2Type constant.WalletType) error {
	wallet1 := new(model.Wallet)
	err := tx.Get(wallet1,
		"SELECT * FROM wallets w WHERE w.account_id = $1 AND w.category = $2 AND w.is_active FOR UPDATE",
		accountId, wallet1Type)
	if err != nil {
		return err
	}

	wallet2 := new(model.Wallet)
	err = tx.Get(wallet2,
		"SELECT * FROM wallets w WHERE w.account_id = $1 AND w.category = $2 AND w.is_active FOR UPDATE",
		accountId, wallet2Type)
	if err != nil {
		return err
	}

	query1 := `
	UPDATE wallets
	SET balance = balance-$1
	WHERE account_id = $2 AND category = $3 AND is_active
	`

	query2 := `
	UPDATE wallets
	SET balance = balance+$1
	WHERE account_id = $2 AND category = $3 AND is_active
	`

	row1, err := tx.Exec(query1, totalPrice, accountId, wallet1Type)
	if err != nil {
		return err
	}
	aff1, _ := row1.RowsAffected()
	if aff1 == 0 {
		return shared.ErrUpdateInactiveWallet
	}

	row2, err := tx.Exec(query2, totalPrice, accountId, wallet2Type)
	if err != nil {
		return err
	}

	aff2, _ := row2.RowsAffected()
	if aff2 == 0 {
		return shared.ErrUpdateInactiveWallet
	}

	return nil
}

func TransferUserTemp(tx *sqlx.Tx, accountId int64, totalPrice decimal.Decimal) error {
	return transferAndRefund(tx, accountId, totalPrice, constant.UserWalletType, constant.TempWalletType)
}

func RefundTempUser(tx *sqlx.Tx, accountId int64, totalPrice decimal.Decimal) error {
	return transferAndRefund(tx, accountId, totalPrice, constant.TempWalletType, constant.UserWalletType)

}

func TransferTempSeller(tx *sqlx.Tx, buyerId, sellerId int64, totalPrice decimal.Decimal) error {
	wallet1 := new(model.Wallet)
	err := tx.Get(wallet1,
		"SELECT * FROM wallets w WHERE w.account_id = $1 AND w.category = $2 AND w.is_active FOR UPDATE",
		buyerId, constant.TempWalletType)
	if err != nil {
		return err
	}

	wallet2 := new(model.Wallet)
	err = tx.Get(wallet2,
		"SELECT * FROM wallets w WHERE w.account_id = $1 AND w.category = $2 AND w.is_active FOR UPDATE",
		sellerId, constant.ShopWalletType)
	if err != nil {
		return err
	}

	query1 := `
	UPDATE wallets
	SET balance = balance-$1
	WHERE account_id = $2 AND category = $3 AND is_active
	`

	query2 := `
	UPDATE wallets
	SET balance = balance+$1
	WHERE account_id = $2 AND category = $3 AND is_active
	`

	row1, err := tx.Exec(query1, totalPrice, buyerId, constant.TempWalletType)
	if err != nil {
		return err
	}
	aff1, _ := row1.RowsAffected()
	if aff1 == 0 {
		return shared.ErrUpdateInactiveWallet
	}

	row2, err := tx.Exec(query2, totalPrice, sellerId, constant.ShopWalletType)
	if err != nil {
		return err
	}
	aff2, _ := row2.RowsAffected()
	if aff2 == 0 {
		return shared.ErrUpdateInactiveWallet
	}

	return nil
}

func (r *walletRepository) FirstShopWalletBalanceBySellerID(ctx context.Context, sellerID int64) (*dto.ShopWalletBalanceResponse, error) {
	balance := new(dto.ShopWalletBalanceResponse)
	qs := `
	SELECT 
		w.balance
	FROM 
		wallets w
	WHERE
		w.category = $1 AND
		w.account_id = $2;
	`

	err := r.db.GetContext(ctx, balance, qs, constant.ShopWalletType, sellerID)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func NewWalletRepository(db *sqlx.DB, tr TransactionRepository) WalletRepository {
	return &walletRepository{
		db: db,
		tr: tr,
	}
}
