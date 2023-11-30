package model

type ReviewMedia struct {
	ID       int64  `db:"id"`
	ReviewID int64  `db:"review_id"`
	ImageUrl string `db:"image_url"`
}
