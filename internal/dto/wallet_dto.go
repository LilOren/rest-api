package dto

import (
	"database/sql"
	"time"

	"github.com/lil-oren/rest/internal/constant"
)

type (
	ActivatePersonalAndTemporaryWalletPayload struct {
		AccountID int64
		Pin       string
	}
	ActivatePersonalAndTemporaryWalletRequestBody struct {
		WalletPin string `json:"wallet_pin" validate:"numeric,len=6"`
	}

	GetPersonalWalletInfoPayload struct {
		UserID int64
	}
	GetPersonalWalletInfoResponse struct {
		IsActive bool    `json:"is_active"`
		Balance  float64 `json:"balance"`
	}
	ListWalletHistoryPayload struct {
		UserID          int64
		StartDate       time.Time
		EndDate         time.Time
		Page            int
		TransactionType constant.ListWalletHistoryQueryValue
	}
	ListWalletHistoryItemFromDB struct {
		Title    string         `db:"title"`
		Amount   float64        `db:"amount"`
		Date     string         `db:"date"`
		IsDebit  bool           `db:"is_debit"`
		ShopName sql.NullString `db:"shop_name"`
		OrderID  sql.NullInt64  `db:"order_id"`
	}
	ListWalletHistoryItem struct {
		Title    string  `json:"title"`
		Amount   float64 `json:"amount"`
		Date     string  `json:"date"`
		IsDebit  bool    `json:"is_debit"`
		ShopName string  `json:"shop_name,omitempty"`
		OrderID  int64   `json:"order_id,omitempty"`
	}

	ListWalletHistoryResponse struct {
		History   []ListWalletHistoryItem `json:"history"`
		Page      int                     `json:"page"`
		TotalPage int                     `json:"total_page"`
	}

	ChangeWalletPinRequestBody struct {
		Password  string `json:"password" validate:"required"`
		WalletPin string `json:"wallet_pin" validate:"required,numeric,len=6"`
	}
	ChangeWalletPinPayload struct {
		UserID    int64
		Password  string
		WalletPin string
	}
)

type (
	SellerWithdrawRequestBody struct {
		Amount float64 `json:"amount" validate:"required,gte=10000"`
	}
	TopUpRequestBody struct {
		Amount float64 `json:"amount" validate:"required,gte=10000"`
	}
	TopUpPayload struct {
		UserID int64
		Amount float64
	}
)

type (
	ShopWalletBalanceResponse struct {
		Balance float64 `json:"balance" db:"balance"`
	}
)
