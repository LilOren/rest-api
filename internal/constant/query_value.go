package constant

type (
	ListWalletHistoryQueryValue string
	CommonQuery                 string
)

const (
	TopupWalletHistoryQueryValue   ListWalletHistoryQueryValue = "topup"
	PaymentWalletHistoryQueryValue ListWalletHistoryQueryValue = "payment"
	RefundWalletHistoryQueryValue  ListWalletHistoryQueryValue = "refund"
	AllWalletHistoryQueryValue     ListWalletHistoryQueryValue = "all"

	StartDateCommonQuery       CommonQuery = "start_date"
	PageCommonQuery            CommonQuery = "page"
	EndDateCommonQuery         CommonQuery = "end_date"
	TransactionTypeCommonQuery CommonQuery = "ttype"
)
