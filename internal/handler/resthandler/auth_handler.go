package resthandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/middleware"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/lil-oren/rest/internal/usecase"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	validate *validator.Validate
	auc      usecase.AuthUsecase
	config   dependency.Config
}

func (h *AuthHandler) register(c *gin.Context) {
	body := new(dto.RegisterUserRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(*body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	if strings.Contains(strings.ToLower(body.Password), strings.ToLower(body.Username)) {
		_ = c.Error(shared.ErrPasswordContainsUsername)
		return
	}

	re := regexp2.MustCompile(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[a-zA-Z]).{8,}$`, regexp2.None)
	passMatch, err := re.MatchString(body.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !passMatch {
		_ = c.Error(shared.ErrPasswordNotMatchRegex)
		return
	}

	ctx := c.Request.Context()
	payload := dto.RegisterUserRequestPayload{
		Username: body.Username,
		Password: body.Password,
		Email:    strings.ToLower(body.Email),
	}
	if err := h.auc.RegisterUser(ctx, payload); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h AuthHandler) login(c *gin.Context) {
	body := new(dto.LoginRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(*body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	payload := dto.LoginRequestPayload(*body)
	resPayload, err := h.auc.Login(ctx, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	shared.SetCookieAfterLogin(c, h.config, resPayload.AccessToken, resPayload.RefreshToken)
	c.Status(http.StatusOK)
}

func (h AuthHandler) refreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	refreshTokenStr, err := c.Cookie(constant.RefreshTokenCookieName)
	if err != nil {
		_ = c.Error(shared.ErrRefreshTokenExpired)
		return
	}

	payload := dto.RefreshTokenPayload{
		RefreshToken: refreshTokenStr,
	}

	resPayload, err := h.auc.RefreshToken(ctx, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	shared.SetCookieAfterRefreshToken(c, h.config, resPayload.AccessToken)

	c.JSON(http.StatusOK, dto.JSONResponse{
		Message: "Successfully refreshed token",
	})
}

func (h AuthHandler) logout(c *gin.Context) {
	ctx := c.Request.Context()

	refreshTokenStr, err := c.Cookie(constant.RefreshTokenCookieName)
	if err != nil {
		_ = c.Error(shared.ErrRefreshTokenExpired)
		return
	}
	_, err = shared.ValidateRefreshToken(refreshTokenStr, h.config)
	if err != nil {
		_ = c.Error(err)
		return
	}

	payload := dto.LogoutPayload{
		RefreshToken: refreshTokenStr,
	}
	if err := h.auc.Logout(ctx, payload); err != nil {
		_ = c.Error(err)
		return
	}

	shared.UnsetCookieAfterLogout(c, h.config)

	c.JSON(http.StatusOK, dto.JSONResponse{
		Message: "Successfully logging out user",
	})
}

func (h AuthHandler) userDetail(c *gin.Context) {
	userId := c.GetInt64(constant.CtxUserId)

	payload := dto.GetUserDetailPayload{
		UserID: userId,
	}

	ctx := c.Request.Context()
	userDetail, err := h.auc.GetUserDetail(ctx, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Data: userDetail,
	})
}

func (h AuthHandler) paymentToken(c *gin.Context) {
	body := new(dto.GetStepUpTokenRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	userID := c.GetInt64(constant.CtxUserId)
	payload := dto.GetStepUpTokenPayload{
		WalletPin: body.WalletPin,
		UserID:    userID,
	}

	ctx := c.Request.Context()
	res, err := h.auc.GetPaymentToken(ctx, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	shared.SetStepUpTokenCookie(c, h.config, res.StepUpToken)
	c.Status(http.StatusOK)
}

func (h AuthHandler) changeEmail(c *gin.Context) {
	body := new(dto.ChangeEmailPayload)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	userID := c.GetInt64(constant.CtxUserId)
	ctx := c.Request.Context()
	payload := dto.ChangeEmailPayload{
		UserID: userID,
		Email:  body.Email,
	}
	if err := h.auc.ChangeEmail(ctx, payload); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h AuthHandler) forgotPassword(c *gin.Context) {
	ctx := c.Request.Context()
	body := new(dto.ForgotPasswordRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	err := h.auc.ForgotPassword(ctx, dto.ForgotPasswordPayload{Email: body.Email})
	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func (h *AuthHandler) resetPassword(c *gin.Context) {
	body := new(dto.ResetPasswordRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(*body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	payload := dto.ResetPasswordPayload{
		ResetCode: body.ResetCode,
		Password:  body.Password,
	}
	if err := h.auc.ResetPassword(ctx, payload); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h AuthHandler) changePasswordRequest(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetInt64(constant.CtxUserId)
	err := h.auc.RequestChangePassword(ctx, userID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func (h AuthHandler) changePassword(c *gin.Context) {
	body := new(dto.ChangePasswordRequestBody)
	if err := c.ShouldBindJSON(body); err != nil {
		_ = c.Error(shared.ErrInvalidBodySchema)
		return
	}

	if err := h.validate.Struct(*body); err != nil {
		_ = c.Error(err.(validator.ValidationErrors))
		return
	}

	ctx := c.Request.Context()
	userID := c.GetInt64(constant.CtxUserId)
	payload := dto.ChangePasswordPayload{
		VerifCode: body.VerifCode,
		Password:  body.Password,
	}
	if err := h.auc.ChangePassword(ctx, payload, userID); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h AuthHandler) googleLogin(c *gin.Context) {
	URL, err := url.Parse(google.Endpoint.AuthURL)
	if err != nil {
		_ = c.Error(err)
		return
	}
	parameters := url.Values{}
	parameters.Add("client_id", h.config.GOauth.ClientID)
	parameters.Add("scope", "https://www.googleapis.com/auth/userinfo.email")
	parameters.Add("redirect_uri", h.config.GOauth.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", "orenlite-state")
	URL.RawQuery = parameters.Encode()
	url := URL.String()

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h AuthHandler) googleLoginCallback(c *gin.Context) {
	oauthConfGl := &oauth2.Config{
		ClientID:     h.config.GOauth.ClientID,
		ClientSecret: h.config.GOauth.ClientSecret,
		RedirectURL:  h.config.GOauth.RedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	state := c.Query("state")
	if state != "orenlite-state" {
		_ = c.Error(shared.ErrInternalServer)
		return
	}

	code := c.Query("code")
	if code == "" {
		reason := c.Query("error_reason")
		if reason == "user_denied" {
			_ = c.Error(shared.ErrInternalServer)
			return
		}
	} else {
		token, err := oauthConfGl.Exchange(c, code)
		if err != nil {
			_ = c.Error(err)
			return
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			_ = c.Error(err)
			return
		}
		defer resp.Body.Close()

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			_ = c.Error(err)
			return
		}
		respStruct := new(dto.GoogleResponse)
		err = json.Unmarshal(response, respStruct)
		if err != nil {
			_ = c.Error(err)
			return
		}
		ctx := c.Request.Context()
		loginResponse, err := h.auc.LoginWithGoogle(ctx, respStruct)
		if err != nil {
			_ = c.Error(err)
			return
		}
		shared.SetCookieAfterLogin(c, h.config, loginResponse.AccessToken, loginResponse.RefreshToken)
		c.Redirect(http.StatusTemporaryRedirect, h.config.GOauth.RedirectFE)
	}
}

func (h AuthHandler) Route(r *gin.Engine) {
	r.
		Group("/auth").
		POST("/register", h.register).
		POST("/login", h.login).
		POST("/refresh-token", h.refreshToken).
		POST("/logout", h.logout).
		GET("/oauth/google", h.googleLogin).
		GET("/oauth/google-callback", h.googleLoginCallback).
		GET("/user", middleware.AllowAuthenticated(h.config), h.userDetail).
		POST("/payment-token", middleware.AllowAuthenticated(h.config), h.paymentToken).
		POST("/change-email", middleware.AllowAuthenticated(h.config), h.changeEmail).
		GET("/hit-auth", middleware.AllowAuthenticated(h.config)).
		POST("/reset-password/request", h.forgotPassword).
		POST("/reset-password", h.resetPassword).
		POST("/change-password/request", middleware.AllowAuthenticated(h.config), h.changePasswordRequest).
		POST("/change-password", middleware.AllowAuthenticated(h.config), h.changePassword)
}

func NewAuthHandler(v *validator.Validate, auc usecase.AuthUsecase, config dependency.Config) *AuthHandler {
	return &AuthHandler{
		validate: v,
		auc:      auc,
		config:   config,
	}
}
