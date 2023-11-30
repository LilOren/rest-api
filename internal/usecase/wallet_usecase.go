package usecase

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

type (
	WalletUsecase interface {
		ActivatePersonalAndTemporaryWallet(ctx context.Context, payload dto.ActivatePersonalAndTemporaryWalletPayload) error
		GetPersonalWalletInfo(ctx context.Context, payload dto.GetPersonalWalletInfoPayload) (*dto.GetPersonalWalletInfoResponse, error)
		SellerWithdrawMoney(ctx context.Context, accountId int64, amount float64) error
		UserTopup(ctx context.Context, payload *dto.TopUpPayload) error
		ListWalletHistory(ctx context.Context, payload dto.ListWalletHistoryPayload) (*dto.ListWalletHistoryResponse, error)
		ChangeWalletPin(ctx context.Context, payload dto.ChangeWalletPinPayload) error
		GetShopWalletBalance(ctx context.Context, sellerID int64) (*dto.ShopWalletBalanceResponse, error)
	}
	walletUsecase struct {
		wr     repository.WalletRepository
		tr     repository.TransactionRepository
		ar     repository.AccountRepository
		config dependency.Config
	}
)

// ChangeWalletPin implements WalletUsecase.
func (uc *walletUsecase) ChangeWalletPin(ctx context.Context, payload dto.ChangeWalletPinPayload) error {
	user, err := uc.ar.FirstById(ctx, payload.UserID)
	if err != nil {
		return err
	}

	if !user.PinHash.Valid {
		return shared.ErrWalletNotActivated
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(payload.Password))
	if err != nil {
		return shared.ErrInvalidPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PinHash.String), []byte(payload.WalletPin))
	if err == nil {
		return shared.ErrSameWalletPin
	}

	pinBytes, _ := bcrypt.GenerateFromPassword([]byte(payload.WalletPin), bcrypt.DefaultCost)

	if err := uc.ar.UpdateWalletPin(ctx, payload.UserID, string(pinBytes)); err != nil {
		return err
	}

	return nil
}

// ListWalletHistory implements WalletUsecase.
func (uc *walletUsecase) ListWalletHistory(ctx context.Context, payload dto.ListWalletHistoryPayload) (*dto.ListWalletHistoryResponse, error) {
	wallet, err := uc.wr.FirstActiveWalletByAccountID(ctx, payload.UserID, constant.UserWalletType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrWalletNotActivated
		}

		return nil, err
	}

	offset := 0
	if payload.Page > 1 {
		offset = payload.Page * 10
	}

	transactions, err := uc.tr.FindTransactionByAccountID(
		ctx,
		wallet.ID,
		payload.StartDate,
		payload.EndDate,
		offset,
	)
	if err != nil {
		return nil, err
	}

	res := dto.ListWalletHistoryResponse{
		History: make([]dto.ListWalletHistoryItem, 0),
		Page:    payload.Page,
	}
	for _, transaction := range transactions {
		temp := dto.ListWalletHistoryItem{
			Title:   transaction.Title,
			Amount:  transaction.Amount,
			Date:    transaction.Date,
			IsDebit: transaction.IsDebit,
		}

		if transaction.OrderID.Valid {
			temp.OrderID = transaction.OrderID.Int64
		}

		if transaction.ShopName.Valid {
			temp.ShopName = transaction.ShopName.String
		}

		res.History = append(res.History, temp)
	}

	count, err := uc.tr.CountTransactionByAccountID(ctx, wallet.ID, payload.StartDate, payload.EndDate)
	if err != nil {
		return nil, err
	}

	res.TotalPage = int(math.Ceil(float64(*count) / 10.0))

	return &res, nil
}

func (uc *walletUsecase) GetPersonalWalletInfo(ctx context.Context, payload dto.GetPersonalWalletInfoPayload) (*dto.GetPersonalWalletInfoResponse, error) {
	wallet, err := uc.wr.FirstActiveWalletByAccountID(ctx, payload.UserID, constant.UserWalletType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrWalletNotActivated
		} else {
			return nil, err
		}
	}

	res := &dto.GetPersonalWalletInfoResponse{
		IsActive: wallet.IsActive,
	}

	res.Balance, _ = wallet.Balance.Float64()

	return res, nil
}

func (uc *walletUsecase) ActivatePersonalAndTemporaryWallet(ctx context.Context, payload dto.ActivatePersonalAndTemporaryWalletPayload) error {
	pinBytes := []byte(payload.Pin)
	pinHashBytes, err := bcrypt.GenerateFromPassword(pinBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	pinHash := string(pinHashBytes)
	err = uc.wr.ActivatePersonalAndTemporaryWallet(ctx, payload.AccountID, pinHash)
	if err != nil {
		return err
	}

	return nil
}

func (uc *walletUsecase) SellerWithdrawMoney(ctx context.Context, accountId int64, amount float64) error {
	walletUser, err := uc.wr.FirstActiveWalletByAccountID(ctx, accountId, constant.UserWalletType)
	if err != nil {
		return err
	}
	if !walletUser.IsActive {
		return shared.ErrUpdateInactiveWallet
	}

	walletShop, err := uc.wr.FirstActiveWalletByAccountID(ctx, accountId, constant.ShopWalletType)
	if err != nil {
		return err
	}

	decAmount := decimal.NewFromFloat(amount)
	if walletShop.Balance.LessThan(decAmount) {
		return shared.ErrInsufficientBalance
	}
	withdraw := &model.Transaction{
		Amount:       decAmount,
		Title:        constant.WithdrawTitle,
		FromWalletID: sql.NullInt64{Int64: walletShop.ID, Valid: true},
		ToWalletID:   walletUser.ID,
	}

	if err = uc.wr.WithdrawShopUser(ctx, accountId, withdraw); err != nil {
		return err
	}
	return nil
}

func (uc *walletUsecase) UserTopup(ctx context.Context, payload *dto.TopUpPayload) error {
	walletUser, err := uc.wr.FirstActiveWalletByAccountID(ctx, payload.UserID, constant.UserWalletType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.ErrWalletNotActivated
		}
		return err
	}
	if !walletUser.IsActive {
		return shared.ErrUpdateInactiveWallet
	}

	topup := &model.Transaction{
		Amount:       decimal.NewFromFloat(payload.Amount),
		Title:        constant.TopUpTitle,
		FromWalletID: sql.NullInt64{},
		ToWalletID:   walletUser.ID,
	}
	if err = uc.wr.Topup(ctx, payload.UserID, topup); err != nil {
		return err
	}
	return nil
}

func (uc *walletUsecase) GetShopWalletBalance(ctx context.Context, sellerID int64) (*dto.ShopWalletBalanceResponse, error) {
	balance, err := uc.wr.FirstShopWalletBalanceBySellerID(ctx, sellerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.ErrWalletNotActivated
		}
		return nil, shared.ErrFindWallet
	}

	return balance, nil
}

func NewWalletUsecase(wr repository.WalletRepository, config dependency.Config,
	tr repository.TransactionRepository, ar repository.AccountRepository) WalletUsecase {
	return &walletUsecase{
		wr:     wr,
		tr:     tr,
		config: config,
		ar:     ar,
	}
}
