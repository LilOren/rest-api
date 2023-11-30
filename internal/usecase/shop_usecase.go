package usecase

import (
	"context"
	"database/sql"
	"math"
	"reflect"
	"strings"

	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
)

type (
	ShopUsecase interface {
		AddShopName(ctx context.Context, payload dto.CreateShopPayload, accountId int) error
		UpdateShopName(ctx context.Context, payload dto.UpdateShopNamePayload, accountId int) error
		UpdateShopAddress(ctx context.Context, payload dto.UpdateShopAddressPayload, accountId int) error
		UpdateShopCourier(ctx context.Context, payload dto.UpdateShopCourierPayload, accountId int) error
		GetCourier(ctx context.Context, accountId int) ([]dto.ShopCourier, error)
		AddProduct(ctx context.Context, payload dto.AddProductPayload, accountId int) error
		GetProductDetail(ctx context.Context, productCode string) (*dto.GetProductDetailResponseBody, error)
		UpdateProduct(ctx context.Context, payload dto.UpdateProductPayload, accountId int) error
		GetAllProduct(ctx context.Context, sellerId int, page int) (*dto.GetAllProductResponseBody, error)
		GetProductDiscount(ctx context.Context, sellerId int, productCode string) (*dto.GetProductDiscountResponseBody, error)
		EditProductDiscount(ctx context.Context, payload dto.UpdateProductDiscountPayload, sellerId int, productCode string) error
		DeleteProduct(ctx context.Context, productCode string, sellerId int64) error
	}
	shopUsecase struct {
		sr  repository.ShopRepository
		aar repository.AccountAddressRepository
		scr repository.ShopCourierRepository
		er  repository.WalletRepository
		pr  repository.ProductRepository
	}
)

func (su *shopUsecase) DeleteProduct(ctx context.Context, productCode string, sellerId int64) error {
	product, err := su.pr.FirstProductByCode(ctx, productCode)
	if err != nil {
		if err == shared.ErrProductNotFound {
			return shared.ErrProductNotFound
		}
		return shared.ErrFindProduct
	}

	if product.SellerID != sellerId {
		return shared.ErrProductNotFromSeller
	}

	err = su.pr.DeleteProduct(ctx, productCode, sellerId)
	if err != nil {
		return shared.ErrDeleteProduct
	}

	return nil
}

func (su *shopUsecase) AddShopName(ctx context.Context, payload dto.CreateShopPayload, accountId int) error {
	_, err := su.sr.FirstShopById(ctx, accountId)
	if err != nil && err != sql.ErrNoRows {
		return shared.ErrFindShop
	}

	if err == nil {
		return shared.ErrAlreadyHaveShop
	}

	_, err = su.er.FirstActiveWalletByAccountID(ctx, int64(accountId), constant.UserWalletType)
	if err != nil && err != sql.ErrNoRows {
		return shared.ErrFindWallet
	}

	if err == sql.ErrNoRows {
		return shared.ErrWalletNotActivated
	}

	err = su.sr.CreateShop(ctx, payload, accountId)
	if err != nil {
		return shared.ErrFailedCreateShop
	}

	err = su.er.ActivateShopWallet(ctx, int64(accountId))
	if err != nil {
		return shared.ErrFailedActivateShopWallet
	}

	return nil
}

func (su *shopUsecase) UpdateShopName(ctx context.Context, payload dto.UpdateShopNamePayload, accountId int) error {
	shop, err := su.sr.FirstShopById(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			return shared.ErrNoShop
		}
		return shared.ErrFindShop
	}

	if shop.Name.String == payload.ShopName {
		return shared.ErrSameShopName
	}

	err = su.sr.UpdateShopName(ctx, payload.ShopName, accountId)
	if err != nil {
		return shared.ErrFailedUpdateShopName
	}

	return nil
}

func (su *shopUsecase) UpdateShopAddress(ctx context.Context, payload dto.UpdateShopAddressPayload, accountId int) error {
	temp := false
	_, err := su.sr.FirstShopById(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			return shared.ErrNoShop
		}
		return shared.ErrFindShop
	}

	address, err := su.aar.FirstShopAddressById(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			return shared.ErrNoAddress
		}
		return shared.ErrFailedGetAddress
	}

	if address.ID == int64(payload.AddressId) {
		return shared.ErrAlreadyDefaultShopAddress
	}

	addressDetailList, err := su.aar.FindDetailsAddressById(ctx, accountId)
	if err != nil {
		return shared.ErrFailedGetAddress
	}

	for _, v := range addressDetailList {
		if v.ID == payload.AddressId {
			temp = true
		}
	}

	if !temp {
		return shared.ErrInvalidAddressId
	}

	err = su.sr.UpdateShopAddress(ctx, payload.AddressId, accountId)
	if err != nil {
		return shared.ErrFailedUpdateShopAddress
	}

	return nil
}

func (su *shopUsecase) UpdateShopCourier(ctx context.Context, payload dto.UpdateShopCourierPayload, accountId int) error {
	courierCond := make([]bool, 0)
	shop, err := su.sr.FirstShopById(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			return shared.ErrNoShop
		}
		return shared.ErrFindShop
	}

	if len(payload) < 3 {
		return shared.ErrCourierNotFull
	}

	if len(payload) > 3 {
		return shared.ErrInvalidCourier
	}

	for _, v := range payload {
		courierCond = append(courierCond, v)
	}

	err = su.sr.UpdateShopCourier(ctx, courierCond, int(shop.ID))
	if err != nil {
		return shared.ErrFailedUpdateShopCourier
	}

	return nil
}

func (su *shopUsecase) UpdateProduct(ctx context.Context, payload dto.UpdateProductPayload, accountId int) error {
	val := reflect.ValueOf(payload.ProductCategoryID)
	categories := make([]interface{}, 0)

	for i := 0; i < val.NumField(); i++ {
		categories = append(categories, val.Field(i).Interface())
	}

	mediaType := make([]string, 0)
	for _, v := range payload.ImageURL {
		switch {
		case strings.Contains(v, constant.MP4VideoType) || strings.Contains(v, constant.MKVVideoType):
			mediaType = append(mediaType, constant.VideoTypeDefault)
		case strings.Contains(v, constant.JPGImageType) || strings.Contains(v, constant.JPEGImageType) || strings.Contains(v, constant.PNGImageType):
			mediaType = append(mediaType, constant.ImageTypeDefault)
		default:
			mediaType = append(mediaType, constant.ImageTypeDefault)
		}
	}

	err := su.pr.UpdateProduct(ctx, payload, categories, mediaType, int64(accountId))
	if err != nil {
		return shared.ErrUpdateProduct
	}
	return nil
}

func (su *shopUsecase) GetCourier(ctx context.Context, accountId int) ([]dto.ShopCourier, error) {
	courierList := make([]dto.ShopCourier, 0)
	courier := dto.ShopCourier{}
	shop, err := su.sr.FirstShopById(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.ErrNoShop
		}
		return nil, shared.ErrFindShop
	}

	shopCourierList, err := su.scr.FindShopCourierByShopId(ctx, shop.ID)
	if err != nil {
		return nil, shared.ErrFailedFindShopCourier
	}

	if len(shopCourierList) < 3 {
		return nil, shared.ErrNoShop
	}

	for _, v := range shopCourierList {
		courier.ID = int(v.ShopCourierID)
		courier.Name = v.CourierName
		courier.ImageURL = v.ImageURL
		courier.Description = v.Description
		courier.IsAvailable = v.IsAvailable
		courierList = append(courierList, courier)
	}

	return courierList, nil
}

func (su *shopUsecase) AddProduct(ctx context.Context, payload dto.AddProductPayload, accountId int) error {
	addProduct := dto.AddProduct{}
	checkDup := make(map[string]int)
	if payload.Weight <= 0 {
		return shared.ErrInvalidWeight
	}

	if payload.ProductCategoryID.Level1 == 0 || payload.ProductCategoryID.Level2 == 0 {
		return shared.ErrNoCategory
	}

	for _, v := range payload.Variants {
		if v.Price < constant.PriceDefault {
			return shared.ErrInvalidPrice
		}
		if v.Stock < constant.StockDefault {
			return shared.ErrInvalidStock
		}
	}

	addProduct.ProductName = payload.ProductName
	addProduct.Description = payload.Description
	addProduct.ImageURL = payload.ImageURL
	addProduct.Weight = payload.Weight
	addProduct.IsVariant = payload.IsVariant
	addProduct.ProductCategoryID = payload.ProductCategoryID

	if !payload.IsVariant {
		payload.VariantDefinitions.VariantGroup1 = &dto.VariantGroup{}
		payload.VariantDefinitions.VariantGroup2 = &dto.VariantGroup{}
		payload.VariantDefinitions.VariantGroup1.Name = constant.ProductVariantDefault
		payload.VariantDefinitions.VariantGroup1.VariantTypes = append(payload.VariantDefinitions.VariantGroup1.VariantTypes, constant.ProductVariantDefault)
		payload.VariantDefinitions.VariantGroup2.Name = constant.ProductVariantDefault
		payload.VariantDefinitions.VariantGroup2.VariantTypes = append(payload.VariantDefinitions.VariantGroup2.VariantTypes, constant.ProductVariantDefault)
	}

	addProduct.VariantDefinitions.VariantGroup1.Name = payload.VariantDefinitions.VariantGroup1.Name
	addProduct.VariantDefinitions.VariantGroup1.VariantTypes = payload.VariantDefinitions.VariantGroup1.VariantTypes

	if payload.VariantDefinitions.VariantGroup2 == nil {
		payload.VariantDefinitions.VariantGroup2 = &dto.VariantGroup{}
		payload.VariantDefinitions.VariantGroup2.Name = constant.ProductVariantDefault
		payload.VariantDefinitions.VariantGroup2.VariantTypes = append(payload.VariantDefinitions.VariantGroup2.VariantTypes, constant.ProductVariantDefault)
	}

	addProduct.VariantDefinitions.VariantGroup2.Name = payload.VariantDefinitions.VariantGroup2.Name
	addProduct.VariantDefinitions.VariantGroup2.VariantTypes = payload.VariantDefinitions.VariantGroup2.VariantTypes

	for idx, v := range addProduct.VariantDefinitions.VariantGroup1.VariantTypes {
		if _, ok := checkDup[strings.ToLower(v)]; ok {
			return shared.ErrDuplicateVariantType
		}
		checkDup[strings.ToLower(v)] = idx
	}

	checkDup = make(map[string]int)
	for idx, v := range addProduct.VariantDefinitions.VariantGroup2.VariantTypes {
		if _, ok := checkDup[strings.ToLower(v)]; ok {
			return shared.ErrDuplicateVariantType
		}
		checkDup[strings.ToLower(v)] = idx
	}

	for idx, v := range payload.Variants {
		if v.VariantType1 == nil {
			payload.Variants[idx].VariantType1 = &payload.VariantDefinitions.VariantGroup1.VariantTypes[0]
		}
		if v.VariantType2 == nil {
			payload.Variants[idx].VariantType2 = &payload.VariantDefinitions.VariantGroup2.VariantTypes[0]
		}
		variant := dto.VariantReq{}
		variant.VariantType1 = *payload.Variants[idx].VariantType1
		variant.VariantType2 = *payload.Variants[idx].VariantType2
		variant.Price = v.Price
		variant.Stock = v.Stock
		addProduct.Variants = append(addProduct.Variants, variant)
	}

	uuid := shared.GenerateUUID()
	mediaType := make([]string, 0)
	productCode := payload.ProductName + "-" + uuid
	for _, v := range payload.ImageURL {
		switch {
		case strings.Contains(v, constant.MP4VideoType) || strings.Contains(v, constant.MKVVideoType):
			mediaType = append(mediaType, constant.VideoTypeDefault)
		case strings.Contains(v, constant.JPGImageType) || strings.Contains(v, constant.JPEGImageType) || strings.Contains(v, constant.PNGImageType):
			mediaType = append(mediaType, constant.ImageTypeDefault)
		default:
			mediaType = append(mediaType, constant.ImageTypeDefault)
		}
	}

	for idx, v := range addProduct.VariantDefinitions.VariantGroup1.VariantTypes {
		if v == constant.ProductVariantDefault {
			break
		}
		if idx == len(addProduct.VariantDefinitions.VariantGroup1.VariantTypes)-1 {
			addProduct.VariantDefinitions.VariantGroup1.VariantTypes = append(addProduct.VariantDefinitions.VariantGroup1.VariantTypes, constant.ProductVariantDefault)
		}
	}

	for idx, v := range addProduct.VariantDefinitions.VariantGroup2.VariantTypes {
		if v == constant.ProductVariantDefault {
			break
		}
		if idx == len(addProduct.VariantDefinitions.VariantGroup2.VariantTypes)-1 {
			addProduct.VariantDefinitions.VariantGroup2.VariantTypes = append(addProduct.VariantDefinitions.VariantGroup2.VariantTypes, constant.ProductVariantDefault)
		}
	}

	err := su.pr.CreateProduct(ctx, addProduct, accountId, productCode, mediaType)
	if err != nil {
		return err
	}

	return nil
}

func (su *shopUsecase) GetProductDetail(ctx context.Context, productCode string) (*dto.GetProductDetailResponseBody, error) {
	productDetail := new(dto.GetProductDetailResponseBody)
	detail, err := su.pr.FindProductDetail(ctx, productCode)
	if err != nil {
		return nil, shared.ErrFindProductDetail
	}

	if len(detail) < 1 {
		return nil, shared.ErrInvalidProductCode
	}

	variant, err := su.pr.FindProductVariants(ctx, int(detail[0].ID))
	if err != nil {
		return nil, shared.ErrFindProductDetail
	}

	productDetail.ID = detail[0].ID
	productDetail.ProductCode = detail[0].ProductCode
	productDetail.ProductName = detail[0].ProductName
	productDetail.Description = detail[0].Description
	productDetail.Weight = detail[0].Weight

	checkDupMedia := make(map[dto.MediaDetail]int)
	checkDupCategory := make(map[dto.CategoryDetail]int)
	checkDupVarGroup := make(map[dto.VariantGroupDetail]int)

	for idx, v := range detail {
		media := dto.MediaDetail{
			MediaID:  int(v.MediaID),
			MediaURL: v.MediaURL,
		}
		category := dto.CategoryDetail{
			CategoryID:   int(v.CategoryID),
			CategoryName: v.CategoryName,
		}
		varGroup1 := dto.VariantGroupDetail{
			VariantGroupID:   int(v.VariantGroup1ID),
			VariantGroupName: v.VariantGroup1Name,
		}
		varGroup2 := dto.VariantGroupDetail{
			VariantGroupID:   int(v.VariantGroup2ID),
			VariantGroupName: v.VariantGroup2Name,
		}
		if _, ok := checkDupMedia[media]; !ok {
			productDetail.Media = append(productDetail.Media, media)
			checkDupMedia[media] = idx
		}
		if _, ok := checkDupCategory[category]; !ok {
			productDetail.Category = append(productDetail.Category, category)
			checkDupCategory[category] = idx
		}
		if _, ok := checkDupVarGroup[varGroup1]; !ok {
			productDetail.VariantDefinition = append(productDetail.VariantDefinition, varGroup1)
			checkDupVarGroup[varGroup1] = idx
		}
		if _, ok := checkDupVarGroup[varGroup2]; !ok {
			productDetail.VariantDefinition = append(productDetail.VariantDefinition, varGroup2)
			checkDupVarGroup[varGroup2] = idx
		}
	}

	for _, v := range variant {
		variants := dto.VariantDetail{
			VariantType1ID:   int(v.VariantType1ID),
			VariantType1Name: v.VariantType1Name,
			VariantType2ID:   int(v.VariantType2ID),
			VariantType2Name: v.VariantType2Name,
			Discount:         v.Discount,
			Price:            v.Price,
			Stock:            v.Stock,
		}
		productDetail.Variants = append(productDetail.Variants, variants)
	}
	return productDetail, nil
}

func (su *shopUsecase) GetAllProduct(ctx context.Context, sellerId int, page int) (*dto.GetAllProductResponseBody, error) {
	res := dto.GetAllProductResponseBody{}

	if page < 0 {
		return nil, shared.ErrInvalidPage
	}
	products, err := su.sr.GetAllProductBySellerId(ctx, sellerId, page)
	if err != nil {
		return nil, shared.ErrFindProduct
	}

	count, err := su.sr.CountAllProductBySellerId(ctx, sellerId)
	if err != nil {
		return nil, shared.ErrFindProduct
	}

	paginationDetail := dto.PaginationDetail{
		Page: page,
	}
	paginationDetail.TotalProduct = *count
	paginationDetail.TotalPage = int(math.Ceil(float64(paginationDetail.TotalProduct) / 10.0))

	res.Products = products
	res.Pagination = paginationDetail

	return &res, nil
}

func (su *shopUsecase) GetProductDiscount(ctx context.Context, sellerId int, productCode string) (*dto.GetProductDiscountResponseBody, error) {
	res := &dto.GetProductDiscountResponseBody{}
	discList := make([]dto.VariantDisc, 0)
	discounts, err := su.sr.FindAllProductDiscountByProductCode(ctx, int64(sellerId), productCode)
	if err != nil {
		return nil, shared.ErrFindProductDiscount
	}

	variantGroup1 := dto.VariantGroup{
		Name: discounts[0].VariantGroup1Name,
	}

	variantGroup2 := dto.VariantGroup{
		Name: discounts[0].VariantGroup2Name,
	}

	checkDupVarType := make(map[string]int)
	for idx, v := range discounts {
		if _, ok := checkDupVarType[v.VariantType1Name]; !ok {
			variantGroup1.VariantTypes = append(variantGroup1.VariantTypes, v.VariantType1Name)
			checkDupVarType[v.VariantType1Name] = idx
		}

		if _, ok := checkDupVarType[v.VariantType2Name]; !ok {
			variantGroup2.VariantTypes = append(variantGroup2.VariantTypes, v.VariantType2Name)
			checkDupVarType[v.VariantType2Name] = idx
		}

		val := dto.VariantDisc{
			VariantType1: v.VariantType1Name,
			VariantType2: v.VariantType2Name,
			Discount:     v.Discount,
		}

		discList = append(discList, val)
	}

	res.VariantDefinition.VariantGroup1 = variantGroup1
	res.VariantDefinition.VariantGroup2 = variantGroup2
	res.Variants = discList

	return res, nil
}

func (su *shopUsecase) EditProductDiscount(ctx context.Context, payload dto.UpdateProductDiscountPayload, sellerId int, productCode string) error {
	err := su.sr.UpdateProductDiscount(ctx, payload, int64(sellerId), productCode)
	if err != nil {
		return shared.ErrUpdateProductDiscount
	}

	return nil
}

func NewShopUsecase(sr repository.ShopRepository, aar repository.AccountAddressRepository, scr repository.ShopCourierRepository, er repository.WalletRepository, pr repository.ProductRepository) ShopUsecase {
	return &shopUsecase{
		sr:  sr,
		aar: aar,
		scr: scr,
		er:  er,
		pr:  pr,
	}
}
