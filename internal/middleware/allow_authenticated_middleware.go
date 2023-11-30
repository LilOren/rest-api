package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/shared"
)

func AllowAuthenticated(config dependency.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("ENV_MODE") == "testing" {
			c.Next()
			return
		}

		accessTokenStr, err := c.Cookie(constant.AccessTokenCookieName)
		if err != nil {
			e := shared.ErrAccessTokenExpired
			c.AbortWithStatusJSON(mapErrorCode[e.Code], e.CreateHTTPErrorMessage())
			return
		}

		token, err := shared.ValidateAccessToken(accessTokenStr, config)
		if err != nil {
			if e, ok := err.(*shared.CustomError); ok {
				c.AbortWithStatusJSON(mapErrorCode[e.Code], e.CreateHTTPErrorMessage())
				return
			}

			e := shared.ErrInvalidToken
			c.AbortWithStatusJSON(mapErrorCode[e.Code], e.CreateHTTPErrorMessage())
			return
		}

		claims, ok := token.Claims.(*shared.AccessJWTClaim)
		if !ok || !token.Valid {
			if err := token.Claims.Valid(); err != nil {
				if e, ok := err.(*shared.CustomError); ok {
					c.AbortWithStatusJSON(mapErrorCode[e.Code], e.CreateHTTPErrorMessage())
					return
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, dto.JSONResponse{
					Message: "internal server error",
				})
				return
			}
			return
		}

		c.Set(constant.CtxUserId, claims.UserId)
		c.Set(constant.CtxIsSeller, claims.IsSeller)

		c.Next()
	}
}
