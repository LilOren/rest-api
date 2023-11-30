package resthandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/middleware"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/lil-oren/rest/internal/usecase"
)

type WalletHandler struct {
	wu     usecase.WalletUsecase
	config dependency.Config
	v      *validator.Validate
}

func (h WalletHandler) activatePersonalWallet(c *gin.Context) {
	body := new(dto.ActivatePersonalAndTemporaryWalletRequestBody)
	if err := c.BindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.v.Struct(body); err != nil {
		e := err.(validator.ValidationErrors)
		_ = c.Error(e)
		return
	}

	ctx := c.Request.Context()
	userID := c.GetInt64(constant.CtxUserId)
	payload := dto.ActivatePersonalAndTemporaryWalletPayload{
		AccountID: int64(userID),
		Pin:       body.WalletPin,
	}

	if err := h.wu.ActivatePersonalAndTemporaryWallet(ctx, payload); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h WalletHandler) getPersonalWalletInfo(c *gin.Context) {
	userID := c.GetInt64(constant.CtxUserId)
	ctx := c.Request.Context()

	p := dto.GetPersonalWalletInfoPayload{
		UserID: userID,
	}

	res, err := h.wu.GetPersonalWalletInfo(ctx, p)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: res,
	})
}

func (h WalletHandler) withdrawMoneySeller(c *gin.Context) {
	req := new(dto.SellerWithdrawRequestBody)
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.v.Struct(req); err != nil {
		e := err.(validator.ValidationErrors)
		_ = c.Error(e)
		return
	}

	ctx := c.Request.Context()
	accountId := c.GetInt64(constant.CtxUserId)
	if err := h.wu.SellerWithdrawMoney(ctx, accountId, req.Amount); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h WalletHandler) topupUser(c *gin.Context) {
	req := new(dto.TopUpRequestBody)
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.v.Struct(req); err != nil {
		e := err.(validator.ValidationErrors)
		_ = c.Error(e)
		return
	}

	ctx := c.Request.Context()
	userID := c.GetInt64(constant.CtxUserId)
	payload := &dto.TopUpPayload{
		UserID: userID,
		Amount: req.Amount,
	}
	if err := h.wu.UserTopup(ctx, payload); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h WalletHandler) listHistory(c *gin.Context) {
	endTime := time.Now().UTC()
	startTime := endTime.AddDate(0, 0, -7)
	var (
		startDateQuery       = c.DefaultQuery(string(constant.StartDateCommonQuery), startTime.Format(constant.DateLayoutISO))
		endDateQuery         = c.DefaultQuery(string(constant.EndDateCommonQuery), endTime.Format(constant.DateLayoutISO))
		transactionTypeQuery = c.DefaultQuery(string(constant.TransactionTypeCommonQuery), string(constant.AllWalletHistoryQueryValue))
		pageQuery            = c.DefaultQuery(string(constant.PageCommonQuery), "1")
	)

	startDate, err := time.Parse(constant.DateLayoutISO, startDateQuery)
	if err != nil {
		_ = c.Error(shared.GenerateErrQueryParamInvalid(string(constant.StartDateCommonQuery)))
		return
	}

	endDate, err := time.Parse(constant.DateLayoutISO, endDateQuery)
	endDate = endDate.Add(time.Hour * 24)
	if err != nil {
		_ = c.Error(shared.GenerateErrQueryParamInvalid(string(constant.EndDateCommonQuery)))
		return
	}

	match, err := shared.WalletHistoryTransactionTypeQueryRegex.MatchString(transactionTypeQuery)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !match {
		_ = c.Error(shared.GenerateErrQueryParamInvalid(string(constant.TransactionTypeCommonQuery)))
		return
	}

	userID := c.GetInt64(constant.CtxUserId)

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		_ = c.Error(err)
		return
	}

	p := dto.ListWalletHistoryPayload{
		UserID:          userID,
		StartDate:       startDate,
		EndDate:         endDate,
		TransactionType: constant.ListWalletHistoryQueryValue(transactionTypeQuery),
		Page:            page,
	}

	ctx := c.Request.Context()
	res, err := h.wu.ListWalletHistory(ctx, p)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: res,
	})
}

func (h WalletHandler) changePin(c *gin.Context) {
	body := new(dto.ChangeWalletPinRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.v.Struct(body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	userID := c.GetInt64(constant.CtxUserId)
	p := dto.ChangeWalletPinPayload{
		UserID:    userID,
		Password:  body.Password,
		WalletPin: body.WalletPin,
	}
	ctx := c.Request.Context()
	if err := h.wu.ChangeWalletPin(ctx, p); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h WalletHandler) getShopWalletBalance(c *gin.Context) {
	sellerID := c.GetInt64(constant.CtxUserId)
	ctx := c.Request.Context()
	res, err := h.wu.GetShopWalletBalance(ctx, sellerID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: res,
	})
}

func (h WalletHandler) Route(r *gin.Engine) {
	r.
		Group("/wallets").
		Use(middleware.AllowAuthenticated(h.config)).
		PUT("/personal/activate", h.activatePersonalWallet).
		GET("/personal/info", h.getPersonalWalletInfo).
		POST("/personal/withdraw", middleware.AllowPayment(h.config), middleware.IsSeller(), h.withdrawMoneySeller).
		POST("/personal/topup", middleware.AllowPayment(h.config), h.topupUser).
		GET("/personal/history", h.listHistory).
		PUT("/change-pin", h.changePin).
		GET("/shop", middleware.IsSeller(), h.getShopWalletBalance)

}

func NewWalletHandler(v *validator.Validate, wu usecase.WalletUsecase, config dependency.Config) WalletHandler {
	return WalletHandler{
		config: config,
		v:      v,
		wu:     wu,
	}
}
