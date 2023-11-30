package constant

type TransactionTitle string

const (
	PaymentOrderTitle TransactionTitle = "PAYMENT-ORDER"
	TopUpTitle        TransactionTitle = "TOPUP"
	WithdrawTitle     TransactionTitle = "WITHDRAW"
	TransferTitle     TransactionTitle = "TRANSFER-ORDER"
	RefundTitle       TransactionTitle = "REFUND-ORDER"
)
