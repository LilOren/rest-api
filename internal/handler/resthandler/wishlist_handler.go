package resthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/middleware"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/lil-oren/rest/internal/usecase"
)

type WishlistHandler struct {
	wu     usecase.WishlistUseCase
	config dependency.Config
	v      *validator.Validate
}

func (h WishlistHandler) addWishlist(c *gin.Context) {
	req := new(dto.WishlistRequestBody)
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.v.Struct(req); err != nil {
		e := err.(validator.ValidationErrors)
		_ = c.Error(e)
		return
	}

	userID := c.GetInt64(constant.CtxUserId)
	ctx := c.Request.Context()
	payload := &dto.WishlistPayload{
		UserID:      userID,
		ProductCode: req.ProductCode,
	}

	if err := h.wu.AddProductToWishlist(ctx, payload); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusCreated)
}

func (h WishlistHandler) deleteWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetInt64(constant.CtxUserId)
	req := new(dto.WishlistRequestBody)
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.v.Struct(req); err != nil {
		e := err.(validator.ValidationErrors)
		_ = c.Error(e)
		return
	}
	payload := &dto.WishlistPayload{
		UserID:      userID,
		ProductCode: req.ProductCode,
	}
	if err := h.wu.DeleteProductFromWishlist(ctx, payload); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func (h WishlistHandler) getAllWishlistUser(c *gin.Context) {
	userID := c.GetInt64(constant.CtxUserId)
	ctx := c.Request.Context()
	params := new(dto.WishlistParams)
	if err := c.ShouldBindQuery(&params); err != nil {
		_ = c.Error(err)
		return
	}
	if err := h.v.Struct(*params); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}
	wishlist, err := h.wu.GetAllProductOfUser(ctx, userID, params)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.JSONResponse{Data: wishlist})
}

func (h WishlistHandler) Route(r *gin.Engine) {
	r.Group("/wishlist", middleware.AllowAuthenticated(h.config)).
		GET("", h.getAllWishlistUser).
		POST("", h.addWishlist).
		DELETE("", h.deleteWishlist)
}

func NewWishlistHandler(wu usecase.WishlistUseCase, config dependency.Config, v *validator.Validate) WishlistHandler {
	return WishlistHandler{
		wu:     wu,
		config: config,
		v:      v,
	}
}
