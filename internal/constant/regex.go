package constant

import "fmt"

const (
	ListProductDistrictQueryRegexPattern = `^(([1-9][0-9]*)(,([1-9][0-9]*))*)?$`
)

var (
	ListWalletHistoryTransactionTypeQueryRegexPattern = fmt.Sprintf(`^(%s|%s|%s|%s)$`,
		TopupWalletHistoryQueryValue,
		PaymentWalletHistoryQueryValue,
		RefundWalletHistoryQueryValue,
		AllWalletHistoryQueryValue,
	)
)
