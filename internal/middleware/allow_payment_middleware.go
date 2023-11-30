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

func AllowPayment(config dependency.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("ENV_MODE") == "testing" {
			c.Next()
			return
		}

		stepUpTokenStr, err := c.Cookie(constant.StepUpTokenCookieName)
		if err != nil {
			e := shared.ErrStepUpTokenExpired
			c.AbortWithStatusJSON(mapErrorCode[e.Code], e.CreateHTTPErrorMessage())
			return
		}

		token, err := shared.ValidateStepUpToken(stepUpTokenStr, config)
		if err != nil {
			if e, ok := err.(*shared.CustomError); ok {
				c.AbortWithStatusJSON(mapErrorCode[e.Code], e.CreateHTTPErrorMessage())
				return
			}

			e := shared.ErrInvalidToken
			c.AbortWithStatusJSON(mapErrorCode[e.Code], e.CreateHTTPErrorMessage())
			return
		}

		_, ok := token.Claims.(*shared.StepUpJWTClaim)
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

		c.Next()
	}
}
