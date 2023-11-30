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

type SellerPageHandler struct {
	spu      usecase.SellerPageUsecase
	config   dependency.Config
	validate *validator.Validate
}

func (h SellerPageHandler) getSellerDetail(c *gin.Context) {
	var (
		pageQuery         = c.DefaultQuery("page", "1")
		sortByQuery       = c.DefaultQuery("sort_by", "most_purchased")
		sortDesc          = c.DefaultQuery("sort_desc", "true")
		searchTerm        = c.DefaultQuery("search_term", "")
		categoryNameQuery = c.DefaultQuery("category_name", "")
	)

	payload := dto.SearchSellerProductPayload{}
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
	payload.CategoryName = categoryNameQuery

	if err := h.validate.Struct(payload); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	shopName := c.Param("name")
	shopDetail, err := h.spu.GetSellerDetail(ctx, shopName, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{Data: shopDetail})
}

func (h SellerPageHandler) Route(r *gin.Engine) {
	r.Group("/shops").
		GET("/:name", h.getSellerDetail)
}

func NewSellerPageHandler(spu usecase.SellerPageUsecase, config dependency.Config, v *validator.Validate) *SellerPageHandler {
	return &SellerPageHandler{
		spu:      spu,
		config:   config,
		validate: v,
	}
}
