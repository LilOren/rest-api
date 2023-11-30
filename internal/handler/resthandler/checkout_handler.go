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

type CheckoutHandler struct {
	cu     usecase.CheckoutUsecase
	v      *validator.Validate
	config dependency.Config
}

func (h CheckoutHandler) summary(c *gin.Context) {
	body := new(dto.CalculateCheckoutSummaryBodyPayload)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	userID := c.GetInt64(constant.CtxUserId)
	payload := dto.CalculateCheckoutSummaryPayload{
		BuyerID:         userID,
		BuyerAddressID:  body.BuyerAddressID,
		OrderDeliveries: body.OrderDeliveries,
	}
	if err := h.v.Struct(payload); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	res, err := h.cu.CalculateCheckoutSummary(ctx, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: res,
	})
}

func (h CheckoutHandler) listCheckoutItem(c *gin.Context) {
	userID := c.GetInt64(constant.CtxUserId)
	ctx := c.Request.Context()
	payload := dto.ListCheckoutItemPayload{
		UserID: userID,
	}

	res, err := h.cu.ListCheckoutItem(ctx, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: res,
	})
}

func (h CheckoutHandler) Route(r *gin.Engine) {
	r.
		Group("/checkouts", middleware.AllowAuthenticated(h.config)).
		GET("", h.listCheckoutItem).
		POST("/summary", h.summary)
}

func NewCheckoutHandler(v *validator.Validate, cu usecase.CheckoutUsecase, config dependency.Config) CheckoutHandler {
	return CheckoutHandler{
		v:      v,
		cu:     cu,
		config: config,
	}
}
