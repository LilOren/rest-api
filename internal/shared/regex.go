package shared

import (
	"github.com/dlclark/regexp2"
	"github.com/lil-oren/rest/internal/constant"
)

var (
	WalletHistoryTransactionTypeQueryRegex = regexp2.MustCompile(constant.ListWalletHistoryTransactionTypeQueryRegexPattern, regexp2.None)
)
