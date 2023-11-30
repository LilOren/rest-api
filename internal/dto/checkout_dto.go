package dto

type (
	OrderDelivery struct {
		ShopID        int64  `json:"shop_id"`
		ShopCourierID *int64 `json:"shop_courier_id,omitempty"`
		PromotionID   int64  `json:"promotion_id,omitempty"`
	}
	CalculateCheckoutSummaryBodyPayload struct {
		OrderDeliveries []OrderDelivery `json:"order_deliveries"`
		BuyerAddressID  int64           `json:"buyer_address_id"`
	}
	CalculateCheckoutSummaryPayload struct {
		BuyerID         int64
		OrderDeliveries []OrderDelivery
		BuyerAddressID  int64
	}

	CalculateCheckoutSummaryOrder struct {
		ShopID            int64   `json:"shop_id"`
		SubTotalProduct   float64 `json:"sub_total_product"`
		SubTotalPromotion float64 `json:"sub_total_promotion"`
		DeliveryCost      float64 `json:"delivery_cost"`
		Subtotal          float64 `json:"subtotal"`
	}
	CalculateCheckoutSummaryResponse struct {
		Orders            []CalculateCheckoutSummaryOrder `json:"orders"`
		TotalShopPrice    float64                         `json:"total_shop_price"`
		TotalProduct      int                             `json:"total_product"`
		TotalDeliveryCost float64                         `json:"total_delivery_cost"`
		ServicePrice      float64                         `json:"service_price"`
		SummaryPrice      float64                         `json:"summary_price"`
	}
	ListCheckoutItemPayload struct {
		UserID int64
	}
	ListCheckoutItem struct {
		Name        string  `json:"name"`
		ImageURL    string  `json:"image_url"`
		Quantity    int     `json:"quantity"`
		TotalWeight int     `json:"total_weight"`
		Price       float64 `json:"price"`
	}
	PromotionDropdown struct {
		PromotionID  int64   `json:"promotion_id"`
		Percentage   float64 `json:"percentage,omitempty"`
		PriceCut     float64 `json:"price_cut,omitempty"`
		MinimumSpend float64 `json:"minimum_spend"`
		IsApplicable bool    `json:"is_applicable"`
	}
	ListCheckout struct {
		ShopID            int64               `json:"shop_id"`
		ShopName          string              `json:"shop_name"`
		ShopCity          string              `json:"shop_city"`
		Items             []ListCheckoutItem  `json:"items"`
		CourierDropdown   []DropdownValue     `json:"courier_dropdown"`
		PromotionDropdown []PromotionDropdown `json:"promotion_dropdown"`
	}
	ListCheckoutItemResponse struct {
		Checkouts         []ListCheckout `json:"checkouts"`
		Balance           float64        `json:"remaining_balance,omitempty"`
		IsWalletActivated bool           `json:"is_wallet_activated"`
		TotalPrice        float64        `json:"total_price"`
	}
)
