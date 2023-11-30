package dto

type (
	DropdownValue struct {
		Label string `json:"label"`
		Value any    `json:"value"`
	}
	ListDistrictPayload struct {
		ProvinceID int64
	}
	ListShopCourierPayload struct {
		ShopID int64
	}
	ListCategoryByParentIDPayload struct {
		ParentCategoryID int64
	}
)
