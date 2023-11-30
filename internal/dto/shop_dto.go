package dto

type (
	CreateShopRequestBody struct {
		ShopName  string `json:"shop_name" validate:"required"`
		AddressId int    `json:"address_id" validate:"required"`
	}
	CreateShopPayload struct {
		ShopName  string
		AddressId int
	}
	CreateShopResponseBody struct {
		ShopName  string `json:"shop_name"`
		AddressId int    `json:"address_id"`
	}
	UpdateShopNameRequestBody struct {
		ShopName string `json:"shop_name" validate:"required"`
	}
	UpdateShopNamePayload struct {
		ShopName string
	}
	UpdateShopNameResponseBody struct {
		ShopName string `json:"shop_name"`
	}
	UpdateShopAddressRequestBody struct {
		AddressId int `json:"address_id" validate:"required"`
	}
	UpdateShopAddressPayload struct {
		AddressId int
	}
	UpdateShopAddressResponseBody struct {
		AddressId int `json:"address_id"`
	}
	UpdateShopCourierRequestBody  map[string]bool
	UpdateShopCourierPayload      map[string]bool
	UpdateShopCourierResponseBody map[string]bool
)

type (
	AddProduct struct {
		ProductName        string
		Description        string
		ImageURL           []string
		Weight             int
		IsVariant          bool
		ProductCategoryID  ProductCategory
		VariantDefinitions VariantDefinitionReq
		Variants           []VariantReq
	}
	AddProductRequestBody struct {
		ProductName        string            `json:"product_name" validate:"required"`
		Description        string            `json:"description" validate:"required"`
		ImageURL           []string          `json:"image_url" validate:"required"`
		Weight             int               `json:"weight" validate:"required"`
		IsVariant          bool              `json:"is_variant"`
		ProductCategoryID  ProductCategory   `json:"product_category_id" validate:"required"`
		VariantDefinitions VariantDefinition `json:"variant_definition"`
		Variants           []Variant         `json:"variants"`
	}
	AddProductPayload struct {
		ProductName        string
		Description        string
		ImageURL           []string
		Weight             int
		IsVariant          bool
		ProductCategoryID  ProductCategory
		VariantDefinitions VariantDefinition
		Variants           []Variant
	}
	ProductCategory struct {
		Level1 int  `json:"level_1" validate:"required"`
		Level2 int  `json:"level_2" validate:"required"`
		Level3 *int `json:"level_3"`
	}
	VariantDefinition struct {
		VariantGroup1 *VariantGroup `json:"variant_group_1"`
		VariantGroup2 *VariantGroup `json:"variant_group_2"`
	}
	VariantGroup struct {
		Name         string   `json:"name"`
		VariantTypes []string `json:"variant_types"`
	}
	Variant struct {
		VariantType1 *string `json:"variant_type1"`
		VariantType2 *string `json:"variant_type2"`
		Price        float64 `json:"price"`
		Stock        int64   `json:"stock"`
	}
	VariantDefinitionReq struct {
		VariantGroup1 VariantGroup `json:"variant_group_1"`
		VariantGroup2 VariantGroup `json:"variant_group_2"`
	}
	VariantGroupReq struct {
		Name         string   `json:"name"`
		VariantTypes []string `json:"variant_types"`
	}
	VariantReq struct {
		VariantType1 string  `json:"variant_type1"`
		VariantType2 string  `json:"variant_type2"`
		Price        float64 `json:"price"`
		Stock        int64   `json:"stock"`
	}
)

type (
	UpdateProduct struct {
		ProductID          int64
		ProductName        string
		Description        string
		ImageURL           []string
		Weight             int
		IsVariant          bool
		ProductCategoryID  ProductCategories
		VariantDefinitions VariantDefinitionReq
		Variants           []VariantReq
	}
	UpdateProductRequestBody struct {
		ProductID          int               `json:"product_id" validate:"required"`
		ProductName        string            `json:"product_name" validate:"required"`
		Description        string            `json:"description" validate:"required"`
		ImageURL           []string          `json:"image_url" validate:"required"`
		Weight             int               `json:"weight" validate:"required"`
		IsVariant          bool              `json:"is_variant"`
		ProductCategoryID  ProductCategories `json:"product_category_id" validate:"required"`
		VariantDefinitions VariantDefinition `json:"variant_definition"`
		Variants           []Variant         `json:"variants"`
	}
	UpdateProductPayload struct {
		ProductID          int
		ProductName        string
		Description        string
		ImageURL           []string
		Weight             int
		IsVariant          bool
		ProductCategoryID  ProductCategories
		VariantDefinitions VariantDefinition
		Variants           []Variant
	}
	ProductCategories struct {
		Level1 int `json:"level_1" validate:"required"`
		Level2 int `json:"level_2" validate:"required"`
		Level3 int `json:"level_3"`
	}
)

type (
	GetAllProduct struct {
		ProductCode  string `db:"product_code"`
		ProductName  string `db:"product_name"`
		ThumbnailURL string `db:"thumbnail_url"`
	}
	GetAllProductResponseBody struct {
		Products   []GetAllProduct  `json:"products"`
		Pagination PaginationDetail `json:"pagination"`
	}
	PaginationDetail struct {
		Page         int `json:"page"`
		TotalPage    int `json:"total_page"`
		TotalProduct int `json:"total_product"`
	}
)

type (
	GetProductDetail struct {
		ID                int64  `db:"id"`
		ProductCode       string `db:"product_code"`
		ProductName       string `db:"product_name"`
		Description       string `db:"description"`
		Weight            int    `db:"weight"`
		MediaID           int64  `db:"media_id"`
		MediaURL          string `db:"media_url"`
		CategoryID        int64  `db:"category_id"`
		CategoryName      string `db:"category_name"`
		VariantGroup1ID   int64  `db:"variant_group1_id"`
		VariantGroup1Name string `db:"variant_group1_name"`
		VariantGroup2ID   int64  `db:"variant_group2_id"`
		VariantGroup2Name string `db:"variant_group2_name"`
	}
	GetProductVariant struct {
		Price            float64 `db:"price"`
		Stock            int     `db:"stock"`
		Discount         float64 `db:"discount"`
		VariantType1ID   int64   `db:"variant_type1_id"`
		VariantType1Name string  `db:"variant_type1_name"`
		VariantType2ID   int64   `db:"variant_type2_id"`
		VariantType2Name string  `db:"variant_type2_name"`
	}
	GetProductDetailResponseBody struct {
		ID                int64                `json:"id"`
		ProductCode       string               `json:"product_code"`
		ProductName       string               `json:"product_name"`
		Description       string               `json:"description"`
		Weight            int                  `json:"weight"`
		Media             []MediaDetail        `json:"media"`
		Category          []CategoryDetail     `json:"category"`
		VariantDefinition []VariantGroupDetail `json:"variant_group"`
		Variants          []VariantDetail      `json:"variant_detail"`
	}
	MediaDetail struct {
		MediaID  int    `json:"media_id"`
		MediaURL string `json:"media_url"`
	}
	CategoryDetail struct {
		CategoryID   int    `json:"category_id"`
		CategoryName string `json:"category_name"`
	}
	VariantGroupDetail struct {
		VariantGroupID   int    `json:"variant_group_id"`
		VariantGroupName string `json:"variant_group_name"`
	}
	VariantDetail struct {
		VariantType1ID   int     `json:"variant_type1_id"`
		VariantType1Name string  `json:"variant_type1_name"`
		VariantType2ID   int     `json:"variant_type2_id"`
		VariantType2Name string  `json:"variant_type2_name"`
		Discount         float64 `json:"discount"`
		Price            float64 `json:"price"`
		Stock            int     `json:"stock"`
	}
)

type (
	GetProductDiscount struct {
		Discount          float64 `db:"discount"`
		VariantGroup1Name string  `db:"variant_group1_name"`
		VariantType1Name  string  `db:"variant_type1_name"`
		VariantGroup2Name string  `db:"variant_group2_name"`
		VariantType2Name  string  `db:"variant_type2_name"`
	}
	GetProductDiscountResponseBody struct {
		VariantDefinition VariantDefinitionReq `json:"variant_group"`
		Variants          []VariantDisc        `json:"variants"`
	}
	VariantDisc struct {
		VariantType1 string  `json:"variant_type1"`
		VariantType2 string  `json:"variant_type2"`
		Discount     float64 `json:"discount"`
	}
)

type (
	UpdateProductDiscountRequestBody struct {
		Variants []VariantDisc `json:"variants"`
	}
	UpdateProductDiscountPayload struct {
		Variants []VariantDisc
	}
)
