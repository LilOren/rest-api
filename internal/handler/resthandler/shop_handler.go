package resthandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/middleware"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/lil-oren/rest/internal/usecase"
)

type ShopHandler struct {
	su       usecase.ShopUsecase
	config   dependency.Config
	validate *validator.Validate
}

func (h ShopHandler) createShop(c *gin.Context) {
	body := new(dto.CreateShopRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(*body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	payload := dto.CreateShopPayload(*body)
	err := h.su.AddShopName(ctx, payload, int(accountId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h ShopHandler) updateShopName(c *gin.Context) {
	body := new(dto.UpdateShopNameRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(*body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	payload := dto.UpdateShopNamePayload(*body)
	err := h.su.UpdateShopName(ctx, payload, int(accountId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h ShopHandler) updateShopAddress(c *gin.Context) {
	body := new(dto.UpdateShopAddressRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(*body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	payload := dto.UpdateShopAddressPayload(*body)
	err := h.su.UpdateShopAddress(ctx, payload, int(accountId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h ShopHandler) updateShopCourier(c *gin.Context) {
	body := new(dto.UpdateShopCourierRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	payload := dto.UpdateShopCourierPayload(*body)
	err := h.su.UpdateShopCourier(ctx, payload, int(accountId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h ShopHandler) getAvailableCourier(c *gin.Context) {
	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	courier, err := h.su.GetCourier(ctx, int(accountId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{Data: courier})
}

func (h ShopHandler) addProduct(c *gin.Context) {
	body := new(dto.AddProductRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(*body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	payload := dto.AddProductPayload(*body)
	err := h.su.AddProduct(ctx, payload, int(accountId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h ShopHandler) getMerchantProductDetail(c *gin.Context) {
	ctx := c.Request.Context()
	productCode := c.Param("code")
	detail, err := h.su.GetProductDetail(ctx, productCode)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{Data: detail})
}

func (h ShopHandler) getProduct(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		_ = c.Error(shared.GenerateErrQueryParamInvalid("page"))
		return
	}

	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	products, err := h.su.GetAllProduct(ctx, int(accountId), page)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{Data: products})
}

func (h ShopHandler) updateProduct(c *gin.Context) {
	body := new(dto.UpdateProductRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(*body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	payload := dto.UpdateProductPayload(*body)
	err := h.su.UpdateProduct(ctx, payload, int(accountId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h ShopHandler) getProductDiscount(c *gin.Context) {
	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	productCode := c.Param("code")
	productDiscount, err := h.su.GetProductDiscount(ctx, int(accountId), productCode)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{Data: productDiscount})
}

func (h ShopHandler) updateProductDiscount(c *gin.Context) {
	body := new(dto.UpdateProductDiscountRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	productCode := c.Param("code")
	payload := dto.UpdateProductDiscountPayload(*body)
	err := h.su.EditProductDiscount(ctx, payload, int(accountId), productCode)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h ShopHandler) deleteProduct(c *gin.Context) {
	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	productCode := c.Param("code")
	err := h.su.DeleteProduct(ctx, productCode, accountId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h ShopHandler) Route(r *gin.Engine) {
	r.
		Group("/merchant", middleware.AllowAuthenticated(h.config)).
		POST("", h.createShop).
		Use(middleware.IsSeller()).
		PUT("/update/name", h.updateShopName).
		PUT("/update/address", h.updateShopAddress).
		PUT("/update/courier", h.updateShopCourier).
		GET("/courier", h.getAvailableCourier).
		POST("/product", h.addProduct).
		GET("/product", h.getProduct).
		PUT("/product", h.updateProduct).
		GET("/product/discount/:code", h.getProductDiscount).
		GET("/product-detail/:code", h.getMerchantProductDetail).
		PUT("/product/discount/:code", h.updateProductDiscount).
		DELETE("/product/:code", h.deleteProduct)
}

func NewShopHandler(su usecase.ShopUsecase, config dependency.Config, v *validator.Validate) *ShopHandler {
	return &ShopHandler{
		su:       su,
		config:   config,
		validate: v,
	}
}
