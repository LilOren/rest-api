package constant

type (
	OrderStatusType   string
)

const (
	NewOrderStatus     OrderStatusType = "NEW"
	ProcessOrderStatus OrderStatusType = "PROCESS"
	DeliverOrderStatus OrderStatusType = "DELIVER"
	ArriveOrderStatus  OrderStatusType = "ARRIVE"
	ReceiveOrderStatus OrderStatusType = "RECEIVE"
	CancelOrderStatus  OrderStatusType = "CANCEL"
)

const (
	DefaultLimitPerPage int = 6
)
