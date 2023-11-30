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

type OrderSellerHandler struct {
	validate *validator.Validate
	ouc      usecase.OrderSellerUsecase
	oc       usecase.OrderUsecase
	cfg      dependency.Config
}

func (h OrderSellerHandler) orderSellerList(c *gin.Context) {
	ctx := c.Request.Context()
	sellerId := c.GetInt64(constant.CtxUserId)
	params := new(dto.OrderSellerParams)
	if err := c.ShouldBindQuery(&params); err != nil {
		_ = c.Error(err)
		return
	}
	if err := h.validate.Struct(*params); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}
	orders, err := h.ouc.GetAllOrderOfSeller(ctx, sellerId, params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	metdata, err := h.ouc.GetAllOrderOfSellerMetadata(ctx, sellerId, params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res := dto.OrderSellerResponse{
		OrdersData: orders,
		TotalData:  metdata.TotalData,
		TotalPage:  metdata.TotalPage,
	}
	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: res,
	})
}

func (h OrderSellerHandler) orderStatusProcess(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	sellerId := c.GetInt64(constant.CtxUserId)
	req := new(dto.OrderSellerStatusRequest)
	req.NewStatus = constant.ProcessOrderStatus
	orderId, err := strconv.Atoi(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err := h.ouc.UpdateOrderStatus(ctx, int64(orderId), req, sellerId); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h OrderSellerHandler) orderStatusDeliver(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	sellerId := c.GetInt64(constant.CtxUserId)
	req := new(dto.OrderSellerStatusRequest)
	req.NewStatus = constant.DeliverOrderStatus
	orderId, err := strconv.Atoi(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}
	if err := h.validate.Struct(*req); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	if err := h.ouc.UpdateOrderStatus(ctx, int64(orderId), req, sellerId); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h OrderSellerHandler) orderStatusArrive(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	sellerId := c.GetInt64(constant.CtxUserId)
	req := new(dto.OrderSellerStatusRequest)
	req.NewStatus = constant.ArriveOrderStatus
	orderId, err := strconv.Atoi(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err := h.ouc.UpdateOrderStatus(ctx, int64(orderId), req, sellerId); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h OrderSellerHandler) orderStatusReject(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	sellerId := c.GetInt64(constant.CtxUserId)
	orderId, err := strconv.Atoi(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err := h.ouc.RejectOrder(ctx, int64(orderId), sellerId); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h OrderSellerHandler) Route(r *gin.Engine) {
	r.
		Group("/orders/seller", middleware.AllowAuthenticated(h.cfg), middleware.IsSeller()).
		GET("", h.orderSellerList).
		PUT("/:id/process", h.orderStatusProcess).
		PUT("/:id/deliver", h.orderStatusDeliver).
		PUT("/:id/arrive", h.orderStatusArrive).
		PUT("/:id/reject", h.orderStatusReject)
}

func NewOrderSellerHandler(v *validator.Validate, ouc usecase.OrderSellerUsecase, cfg dependency.Config, oc usecase.OrderUsecase) *OrderSellerHandler {
	return &OrderSellerHandler{
		validate: v,
		ouc:      ouc,
		cfg:      cfg,
		oc:       oc,
	}
}
