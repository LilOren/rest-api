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

type OrderHandler struct {
	ou       usecase.OrderUsecase
	config   dependency.Config
	validate *validator.Validate
}

func (h OrderHandler) createOrder(c *gin.Context) {
	body := new(dto.CreateOrderRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	payload := dto.CreateOrderRequestPayload(*body)
	err := h.ou.CreateOrder(ctx, int(accountId), payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h OrderHandler) orderList(c *gin.Context) {
	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	params := new(dto.OrderParams)
	if err := c.ShouldBindQuery(&params); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.validate.Struct(*params); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	orders, err := h.ou.GetAllOrder(ctx, accountId, params)
	if err != nil {
		_ = c.Error(err)
		return
	}
	totalPage, err := h.ou.GetAllOrderMetadata(ctx, accountId, params)
	if err != nil {
		_ = c.Error(err)
		return
	}
	res := dto.OrderBuyerResponse{
		Orders:     orders,
		Pagination: *totalPage,
	}
	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: res,
	})
}

func (h *OrderHandler) orderStatusReceive(c *gin.Context) {
	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	id := c.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.ou.ReceiveOrder(ctx, int64(orderId), accountId); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *OrderHandler) orderStatusCancel(c *gin.Context) {
	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	id := c.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.ou.CancelOrder(ctx, int64(orderId), accountId); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *OrderHandler) Route(r *gin.Engine) {
	r.
		Group("/orders", middleware.AllowAuthenticated(h.config)).
		POST("", middleware.AllowPayment(h.config), h.createOrder).
		GET("", h.orderList).
		PUT("/:id/receive", h.orderStatusReceive).
		PUT("/:id/cancel", h.orderStatusCancel)
}

func NewOrderHandler(ou usecase.OrderUsecase, config dependency.Config, v *validator.Validate) *OrderHandler {
	return &OrderHandler{
		ou:       ou,
		config:   config,
		validate: v,
	}
}
