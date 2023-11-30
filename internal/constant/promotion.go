package constant

type (
	PromotionStatusType string
)

const (
	OngoingPromotionStatus PromotionStatusType = "ONGOING"
	ComingPromotionStatus  PromotionStatusType = "COMING"
	EndedPromotionStatus   PromotionStatusType = "ENDED"
)
