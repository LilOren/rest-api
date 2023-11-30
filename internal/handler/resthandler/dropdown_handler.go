package resthandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/lil-oren/rest/internal/usecase"
)

type DropdownHandler struct {
	validate *validator.Validate
	dr       usecase.DropdownUsecase
	cfg      dependency.Config
}

func (h DropdownHandler) listProvinces(c *gin.Context) {
	ctx := c.Request.Context()
	provinces, err := h.dr.ListProvince(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: provinces,
	})
}

func (h DropdownHandler) listDistrictsByProvinceID(c *gin.Context) {
	provinceIdStr := c.Params.ByName("province_id")
	provinceId, err := strconv.Atoi(provinceIdStr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	ctx := c.Request.Context()
	payload := dto.ListDistrictPayload{
		ProvinceID: int64(provinceId),
	}
	districts, err := h.dr.ListDistrict(ctx, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: districts,
	})
}

func (h DropdownHandler) listShopCourier(c *gin.Context) {
	shopIDStr := c.Query("shop_id")
	if shopIDStr == "" {
		_ = c.Error(shared.GenerateErrQueryParamRequired("shop_id"))
		return
	}

	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	ctx := c.Request.Context()
	payload := dto.ListShopCourierPayload{
		ShopID: int64(shopID),
	}

	couriers, err := h.dr.ListShopCourier(ctx, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: couriers,
	})
}

func (h DropdownHandler) listTopCategory(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := h.dr.ListTopCategory(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: res,
	})
}

func (h DropdownHandler) listChildCategory(c *gin.Context) {
	parentIDQuery := c.Query("parent_id")
	if parentIDQuery == "" {
		_ = c.Error(shared.GenerateErrQueryParamRequired("parent_id"))
		return
	}

	parentID, err := strconv.Atoi(parentIDQuery)
	if err != nil {
		_ = c.Error(shared.GenerateErrQueryParamInvalid("parent_id"))
		return
	}

	p := dto.ListCategoryByParentIDPayload{
		ParentCategoryID: int64(parentID),
	}
	ctx := c.Request.Context()
	res, err := h.dr.ListCategoryByParentID(ctx, p)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: res,
	})
}

func (h DropdownHandler) Route(r *gin.Engine) {
	r.
		Group("/dropdowns").
		GET("/location-unit/provinces", h.listProvinces).
		GET("/location-unit/provinces/:province_id/districts", h.listDistrictsByProvinceID).
		GET("/products/top-categories", h.listTopCategory).
		GET("/products/child-category", h.listChildCategory).
		GET("/checkouts/couriers", h.listShopCourier)
}

func NewDropdownHandler(v *validator.Validate, dr usecase.DropdownUsecase, cfg dependency.Config) *DropdownHandler {
	return &DropdownHandler{
		validate: v,
		dr:       dr,
		cfg:      cfg,
	}
}
