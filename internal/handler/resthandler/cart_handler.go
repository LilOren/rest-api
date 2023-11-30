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

type CartHandler struct {
	validate *validator.Validate
	cu       usecase.CartUsecase
	cfg      dependency.Config
}

func (h CartHandler) getProduct(c *gin.Context) {
	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	res, err := h.cu.GetCartPageAllProducts(ctx, accountId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.JSONResponse{Data: res})
}

func (h CartHandler) addToCart(c *gin.Context) {
	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	product := new(dto.AddToCartRequestBody)
	if err := c.ShouldBindJSON(product); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}
	if err := h.validate.Struct(*product); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}
	productPayload := &dto.AddToCartRequestPayload{
		ProductVariantID: product.ProductVariantID,
		SellerID:         product.SellerID,
		Quantity:         product.Quantity,
	}
	err := h.cu.AddToCart(ctx, productPayload, accountId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h CartHandler) updateQuantityItem(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	accountId := c.GetInt64(constant.CtxUserId)
	cartId, err := strconv.Atoi(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	req := new(dto.UpdateQuantityRequestBody)
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}
	if err := h.validate.Struct(*req); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}
	err = h.cu.UpdateQuantityItem(ctx, int64(cartId), int(req.Quantity))
	if err != nil {
		_ = c.Error(err)
		return
	}
	prices, err := h.cu.GetTotalPriceChecked(ctx, accountId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.JSONResponse{Data: prices})
}

func (h CartHandler) deleteItem(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	cartId, err := strconv.Atoi(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err = h.cu.DeleteItem(ctx, int64(cartId)); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func (h CartHandler) checkItems(c *gin.Context) {
	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	req := new(dto.IsCheckedCartRequestBody)
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}
	if err := h.validate.Struct(*req); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}
	if err := h.cu.UpdateIsCheckCart(ctx, req.IsCheckCarts); err != nil {
		_ = c.Error(err)
		return
	}
	res, err := h.cu.GetTotalPriceChecked(ctx, accountId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.JSONResponse{Data: res})
}

func (h CartHandler) Route(r *gin.Engine) {
	r.
		Group("/carts", middleware.AllowAuthenticated(h.cfg)).
		GET("", h.getProduct).
		POST("", h.addToCart).
		PUT("/:id", h.updateQuantityItem).
		DELETE("/:id", h.deleteItem).
		PUT("/check-items", h.checkItems)
}

func NewCartHandler(v *validator.Validate, cu usecase.CartUsecase, cfg dependency.Config) *CartHandler {
	return &CartHandler{
		validate: v,
		cu:       cu,
		cfg:      cfg,
	}
}
