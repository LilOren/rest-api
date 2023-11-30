package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/shared"
)

func GetUserID(config dependency.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenStr, err := c.Cookie(constant.AccessTokenCookieName)
		if err != nil {
			c.Next()
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
		c.Next()
	}
}
