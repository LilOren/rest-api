package dto

type (
	ProductPageResponse struct {
		Product         *ProductPageProductDetail   `json:"product"`
		Shop            *ProductPageShop            `json:"shop"`
		ProductVariants []ProductPageProductVariant `json:"product_variant"`
		ProductMedias   []ProductPageProductMedia   `json:"product_media"`
		VariantGroup1   ProductPageVariantsResponse `json:"variant_group1"`
		VariantGroup2   ProductPageVariantsResponse `json:"variant_group2"`
		HighPrice       float64                     `json:"high_price"`
		LowPrice        float64                     `json:"low_price"`
		IsVariant       bool                        `json:"is_variant"`
		WishlistCtr     int                         `json:"wishlist_count"`
		IsInWishlist    bool                        `json:"is_in_wishlist"`
		Rating          float64                     `json:"rating"`
		ReviewCount     int                         `json:"rating_count"`
		TotalSold       int                         `json:"total_sold"`
	}
)

type (
	ProductPageProductDetail struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Weight      int    `json:"weight"`
	}
	ProductPageShop struct {
		ID                int64  `db:"id" json:"id"`
		Name              string `db:"name" json:"name"`
		ProfilePictureURL string `db:"profile_picture_url" json:"profile_picture_url"`
		Location          string `db:"location" json:"location"`
	}
	ProductPageProductVariant struct {
		ID              int64   `json:"id"`
		Price           float64 `json:"price"`
		DiscountedPrice float64 `json:"discounted_price"`
		Stock           uint32  `json:"stock"`
		Discount        float32 `json:"discount"`
		VariantType1ID  int64   `json:"variant_type1_id"`
		VariantType2ID  int64   `json:"variant_type2_id"`
	}
	ProductPageProductMedia struct {
		MediaUrl  string `json:"media_url"`
		MediaType string `json:"mediat_type"`
	}
	ProductPageVariants struct {
		GroupName string `db:"group_name"`
		TypeID    int64  `db:"type_id"`
		TypeName  string `db:"type_name"`
	}
	ProductPageVariantType struct {
		TypeID   int64  `json:"type_id"`
		TypeName string `json:"type_name"`
	}
	ProductPageVariantsResponse struct {
		GroupName    string                   `json:"group_name"`
		VariantTypes []ProductPageVariantType `json:"variant_types"`
	}
)
