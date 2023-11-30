package resthandler

import (
	"net/http"
	"strconv"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/middleware"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/lil-oren/rest/internal/usecase"
)

type ProductPageHandler struct {
	validate *validator.Validate
	cfg      dependency.Config
	puc      usecase.ProductPageUsecase
	dc       usecase.DiscoveryUsecase
}

func (h ProductPageHandler) productDetail(c *gin.Context) {
	productCode := c.Param("product_code")

	ctx := c.Request.Context()
	userID := c.GetInt64(constant.CtxUserId)
	res, err := h.puc.GetProductDetail(ctx, productCode, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.JSONResponse{Data: res})
}

func (h ProductPageHandler) listProduct(c *gin.Context) {
	var (
		pageQuery       = c.DefaultQuery("page", "1")
		sortByQuery     = c.DefaultQuery("sort_by", "created_at")
		sortDesc        = c.DefaultQuery("sort_desc", "true")
		districts       = c.DefaultQuery("districts", "")
		searchTerm      = c.DefaultQuery("search_term", "")
		categoryIdQuery = c.DefaultQuery("category", "")
		maxPriceQuery   = c.DefaultQuery("max_price", "0")
		minPriceQuery   = c.DefaultQuery("min_price", "0")
	)

	payload := dto.SearchProductPayload{
		MaxPrice: 0,
		MinPrice: 0,
	}
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		_ = c.Error(shared.GenerateErrQueryParamInvalid("page"))
		return
	}

	payload.Page = page
	payload.SortBy = sortByQuery
	sd, err := strconv.ParseBool(sortDesc)
	if err != nil {
		_ = c.Error(err)
		return
	}

	payload.SortDesc = sd
	payload.SearchTerm = searchTerm

	districtsRegexPattern := regexp2.MustCompile(constant.ListProductDistrictQueryRegexPattern, regexp2.None)
	match, err := districtsRegexPattern.MatchString(districts)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !match {
		_ = c.Error(shared.GenerateErrQueryParamInvalid("districts"))
		return
	}

	payload.DistrictIDs = districts

	if categoryIdQuery != "" {
		categoryID, err := strconv.Atoi(categoryIdQuery)
		if err != nil {
			_ = c.Error(shared.GenerateErrQueryParamInvalid("category_id"))
			return
		}

		payload.CategoryID = int64(categoryID)
	}

	minPrice, err := strconv.ParseFloat(minPriceQuery, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}

	maxPrice, err := strconv.ParseFloat(maxPriceQuery, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}

	payload.MaxPrice = maxPrice
	payload.MinPrice = minPrice

	if err := h.validate.Struct(payload); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	res, err := h.dc.SearchProduct(ctx, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: res,
	})
}

func (h ProductPageHandler) Route(r *gin.Engine) {
	r.
		Group("/products").
		GET("", h.listProduct).
		GET("/:product_code", middleware.GetUserID(h.cfg), h.productDetail)
}

func NewProductPageHandler(v *validator.Validate, puc usecase.ProductPageUsecase, dc usecase.DiscoveryUsecase, cfg dependency.Config) *ProductPageHandler {
	return &ProductPageHandler{
		validate: v,
		cfg:      cfg,
		puc:      puc,
		dc:       dc,
	}
}
