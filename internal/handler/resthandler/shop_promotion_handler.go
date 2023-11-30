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

type ShopPromotionHandler struct {
	puc usecase.PromotionUsecase
	cfg dependency.Config
	v   *validator.Validate
}

func (h ShopPromotionHandler) getAllPromoFromShop(c *gin.Context) {
	params := new(dto.PromotionByShopParams)
	if err := c.ShouldBindQuery(&params); err != nil {
		_ = c.Error(err)
		return
	}
	if err := h.v.Struct(*params); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}
	sellerID := c.GetInt64(constant.CtxUserId)
	ctx := c.Request.Context()
	res, err := h.puc.GetAllPromoFromShop(ctx, sellerID, params)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.JSONResponse{Data: res})
}

func (h ShopPromotionHandler) addShopPromotion(c *gin.Context) {
	req := new(dto.UpsertShopPromotionRequestBody)
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.v.Struct(req); err != nil {
		e := err.(validator.ValidationErrors)
		_ = c.Error(e)
		return
	}

	sellerID := c.GetInt64(constant.CtxUserId)
	ctx := c.Request.Context()
	payload := &dto.UpsertShopPromotionPayload{
		SellerID:     sellerID,
		Name:         req.Name,
		ExactPrice:   req.ExactPrice,
		Percentage:   req.Percentage,
		MinimumSpend: req.MinimumSpend,
		Quota:        req.Quota,
		StartedAt:    req.StartedAt,
		ExpiredAt:    req.ExpiredAt,
	}

	if err := h.puc.AddShopPromotion(ctx, payload); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusCreated)
}

func (h ShopPromotionHandler) updateShopPromotion(c *gin.Context) {
	req := new(dto.UpsertShopPromotionRequestBody)
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.v.Struct(req); err != nil {
		e := err.(validator.ValidationErrors)
		_ = c.Error(e)
		return
	}
	id := c.Param("id")
	promoID, err := strconv.Atoi(id)
	if err != nil {
		c.Error(err)
		return
	}

	sellerID := c.GetInt64(constant.CtxUserId)
	ctx := c.Request.Context()
	payload := &dto.UpsertShopPromotionPayload{
		SellerID:     sellerID,
		Name:         req.Name,
		ExactPrice:   req.ExactPrice,
		Percentage:   req.Percentage,
		MinimumSpend: req.MinimumSpend,
		Quota:        req.Quota,
		StartedAt:    req.StartedAt,
		ExpiredAt:    req.ExpiredAt,
	}

	if err := h.puc.UpdateShopPromotion(ctx, payload, int64(promoID)); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func (h ShopPromotionHandler) duplicateShopPromotion(c *gin.Context) {
	id := c.Param("id")
	promoID, err := strconv.Atoi(id)
	if err != nil {
		c.Error(err)
		return
	}

	sellerID := c.GetInt64(constant.CtxUserId)
	ctx := c.Request.Context()
	if err := h.puc.DuplicateShopPromotion(ctx, int64(promoID), sellerID); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func (h ShopPromotionHandler) deleteShopPromotion(c *gin.Context) {
	id := c.Param("id")
	promoID, err := strconv.Atoi(id)
	if err != nil {
		c.Error(err)
		return
	}

	sellerID := c.GetInt64(constant.CtxUserId)
	ctx := c.Request.Context()
	if err := h.puc.DeleteShopPromotion(ctx, int64(promoID), sellerID); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func (h ShopPromotionHandler) Route(r *gin.Engine) {
	r.Group("/shop-promotions", middleware.AllowAuthenticated(h.cfg), middleware.IsSeller()).
		GET("", h.getAllPromoFromShop).
		POST("", h.addShopPromotion).
		PUT("/:id", h.updateShopPromotion).
		POST("/:id/duplicate", h.duplicateShopPromotion).
		DELETE("/:id", h.deleteShopPromotion)
}

func NewShopPromotionHandler(puc usecase.PromotionUsecase, cfg dependency.Config, v *validator.Validate) ShopPromotionHandler {
	return ShopPromotionHandler{
		puc: puc,
		cfg: cfg,
		v:   v,
	}
}
