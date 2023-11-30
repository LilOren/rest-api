package model

type ShopCourier struct {
	ID          int64 `db:"id"`
	ShopId      int64 `db:"shop_id"`
	CourierId   int64 `db:"courier_id"`
	IsAvailable bool  `db:"is_available"`
}
