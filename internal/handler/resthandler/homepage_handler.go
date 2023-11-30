package resthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/middleware"
	"github.com/lil-oren/rest/internal/usecase"
)

type HomePageHandler struct {
	validate *validator.Validate
	huc      usecase.HomepageUsecase
	cfg      dependency.Config
}

func (h HomePageHandler) homePageProduct(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := h.huc.GetRecommendedProducts(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{Data: res})
}

func (h HomePageHandler) homePageCart(c *gin.Context) {
	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	res, err := h.huc.GetCartForHome(ctx, accountId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.JSONResponse{Data: res})
}

func (h HomePageHandler) listCategories(c *gin.Context) {
	ctx := c.Request.Context()
	cats, err := h.huc.GetTopCategories(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: cats,
	})
}

func (h HomePageHandler) Route(r *gin.Engine) {
	r.Group("/home-page").
		GET("/recommended-products", h.homePageProduct).
		GET("/carts", middleware.AllowAuthenticated(h.cfg), h.homePageCart).
		GET("/categories", h.listCategories)

}

func NewHomePageHandler(
	v *validator.Validate,
	huc usecase.HomepageUsecase,
	cfg dependency.Config,
) *HomePageHandler {
	return &HomePageHandler{
		validate: v,
		huc:      huc,
		cfg:      cfg,
	}
}
