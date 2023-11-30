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

type ProfileHandler struct {
	pu       usecase.ProfileUsecase
	validate *validator.Validate
	config   dependency.Config
}

func (h ProfileHandler) addAddress(c *gin.Context) {
	body := new(dto.AddAddressRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(*body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	accountID := c.GetInt64(constant.CtxUserId)
	payload := dto.AddAddressPayload(*body)
	err := h.pu.AddAddress(ctx, payload, int(accountID))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h ProfileHandler) changeDefaultAddress(c *gin.Context) {
	body := new(dto.AccountDetailsAddressBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(*body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	accountID := c.GetInt64(constant.CtxUserId)
	payload := dto.AccountDetailsAddressPayload(*body)
	err := h.pu.ChangeDefaultAddress(ctx, int(accountID), payload.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h ProfileHandler) getAccountAddress(c *gin.Context) {
	ctx := c.Request.Context()
	accountID := c.GetInt64(constant.CtxUserId)
	res, err := h.pu.GetAddressDetailsByAccountId(ctx, int(accountID))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: res,
	})
}

func (h ProfileHandler) updatePicture(c *gin.Context) {
	body := new(dto.UploadProfilePictureRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.validate.Struct(body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	userID := c.GetInt64(constant.CtxUserId)
	if err := h.pu.UploadProfilePicture(ctx, userID, body.ImageURL); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h ProfileHandler) updateAddress(c *gin.Context) {
	body := new(dto.UpdateAddressByIDRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	addressIDStr := c.Params.ByName("id")
	addressID, err := strconv.Atoi(addressIDStr)
	if err != nil {
		_ = c.Error(shared.GenerateErrQueryParamInvalid("id"))
		return
	}

	ctx := c.Request.Context()
	userID := c.GetInt64(constant.CtxUserId)
	p := dto.UpdateAddressByIDPayload{
		UserID:              userID,
		AddressID:           int64(addressID),
		ReceiverName:        body.ReceiverName,
		ReceiverPhoneNumber: body.ReceiverPhoneNumber,
		Address:             body.Address,
		PostalCode:          body.PostalCode,
		ProvinceID:          body.ProvinceID,
		DistrictID:          body.DistrictID,
	}
	if err := h.pu.UpdateAddress(ctx, p); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h ProfileHandler) getAccountAddressDetail(c *gin.Context) {
	addressIDStr := c.Params.ByName("id")
	addressID, err := strconv.Atoi(addressIDStr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userID := c.GetInt64(constant.CtxUserId)
	p := dto.GetAddressByIDPayload{
		AccountID: userID,
		AddressID: int64(addressID),
	}

	ctx := c.Request.Context()
	res, err := h.pu.GetAddressDetailByID(ctx, p)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: res,
	})
}

func (h ProfileHandler) Route(r *gin.Engine) {
	r.
		Group("/profile").
		Use(middleware.AllowAuthenticated(h.config)).
		GET("/addresses", h.getAccountAddress).
		GET("/addresses/:id", h.getAccountAddressDetail).
		POST("/addresses", h.addAddress).
		PUT("/addresses/:id", h.updateAddress).
		PUT("/addresses/change-default", h.changeDefaultAddress).
		PUT("/picture", h.updatePicture)
}

func NewProfileHandler(
	pu usecase.ProfileUsecase,
	config dependency.Config,
	v *validator.Validate,
) ProfileHandler {
	return ProfileHandler{
		pu:       pu,
		config:   config,
		validate: v,
	}
}
