package shared

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
)

func SetCookieAfterLogin(c *gin.Context, config dependency.Config, accessToken, refreshToken string) {
	accessTokenCookieExp := int(config.Jwt.AccessTokenExpiration) * 60
	refreshTokenCookieExp := int(config.Jwt.RefreshTokenExpiration) * 60

	if config.App.OriginDomain == "localhost" {
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(constant.AccessTokenCookieName, accessToken, accessTokenCookieExp, "/", config.App.OriginDomain, false, true)
		c.SetCookie(constant.RefreshTokenCookieName, refreshToken, refreshTokenCookieExp, "/", config.App.OriginDomain, false, true)
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.AccessTokenCookieName, accessToken, accessTokenCookieExp, "/", config.App.OriginDomain, true, true)
	c.SetCookie(constant.RefreshTokenCookieName, refreshToken, refreshTokenCookieExp, "/", config.App.OriginDomain, true, true)
}

func UnsetCookieAfterLogout(c *gin.Context, config dependency.Config) {

	if config.App.OriginDomain == "localhost" {
		c.SetCookie(constant.AccessTokenCookieName, "", -1, "/", config.App.OriginDomain, false, true)
		c.SetCookie(constant.RefreshTokenCookieName, "", -1, "/", config.App.OriginDomain, false, true)
		return
	}

	c.SetCookie(constant.AccessTokenCookieName, "", -1, "/", config.App.OriginDomain, true, true)
	c.SetCookie(constant.RefreshTokenCookieName, "", -1, "/", config.App.OriginDomain, true, true)
}

func SetCookieAfterRefreshToken(c *gin.Context, config dependency.Config, accessToken string) {
	accessTokenCookieExp := int(config.Jwt.AccessTokenExpiration) * 60

	if config.App.OriginDomain == "localhost" {
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(constant.AccessTokenCookieName, accessToken, accessTokenCookieExp, "/", config.App.OriginDomain, false, true)
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.AccessTokenCookieName, accessToken, accessTokenCookieExp, "/", config.App.OriginDomain, true, true)
}

func SetStepUpTokenCookie(c *gin.Context, config dependency.Config, stepUpToken string) {
	stepUpTokenCookieExp := int(config.Jwt.StepUpTokenExpiration) * 60

	if config.App.OriginDomain == "localhost" {
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(constant.StepUpTokenCookieName, stepUpToken, stepUpTokenCookieExp, "/", config.App.OriginDomain, false, true)
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.StepUpTokenCookieName, stepUpToken, stepUpTokenCookieExp, "/", config.App.OriginDomain, true, true)
}
