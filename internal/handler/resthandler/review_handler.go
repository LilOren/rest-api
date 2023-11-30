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

type ReviewHandler struct {
	ruc    usecase.ReviewUsecase
	config dependency.Config
	v      *validator.Validate
}

func (h ReviewHandler) addReviewOfProduct(c *gin.Context) {
	req := new(dto.AddReviewRequestBody)
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
	payload := &dto.AddReviewPayload{
		AccountID:   userID,
		Rating:      req.Rating,
		Comment:     req.Comment,
		ProductCode: req.ProductCode,
		ImageUrls:   req.ImageUrls,
	}

	if err := h.ruc.AddReviewOfProduct(ctx, payload); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusCreated)
}

func (h ReviewHandler) getReviewOfProduct(c *gin.Context) {
	productCode := c.Param("product_code")
	params := new(dto.GetAllReviewParams)
	if err := c.ShouldBindQuery(&params); err != nil {
		_ = c.Error(err)
		return
	}
	if err := h.v.Struct(*params); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}
	ctx := c.Request.Context()
	res, err := h.ruc.GetAllReviewOfProduct(ctx, productCode, params)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.JSONResponse{Data: res})
}

func (h ReviewHandler) Route(r *gin.Engine) {
	r.Group("/reviews").
		GET("/:product_code", h.getReviewOfProduct).
		Use(middleware.AllowAuthenticated(h.config)).
		POST("", h.addReviewOfProduct)
}

func NewReviewHandler(ruc usecase.ReviewUsecase, config dependency.Config, v *validator.Validate) ReviewHandler {
	return ReviewHandler{
		ruc:    ruc,
		config: config,
		v:      v,
	}
}
