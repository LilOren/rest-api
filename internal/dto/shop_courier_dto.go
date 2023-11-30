package dto

type (
	ShopCourier struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		ImageURL    string `json:"image_url"`
		Description string `json:"description"`
		IsAvailable bool   `json:"is_available"`
	}
	ShopCourierDetailsResponse struct {
		ShopCourierID int64  `db:"shop_courier_id"`
		CourierName   string `db:"courier_name"`
		ImageURL      string `db:"image_url"`
		Description   string `db:"description"`
		IsAvailable   bool   `db:"is_available"`
	}
)
