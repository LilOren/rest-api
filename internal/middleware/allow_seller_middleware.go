package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/shared"
)

func IsSeller() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.GetBool(constant.CtxIsSeller) {
			e := shared.ErrUnauthorizedUser
			c.AbortWithStatusJSON(mapErrorCode[e.Code], e.CreateHTTPErrorMessage())
			return
		}

		c.Next()
	}
}
