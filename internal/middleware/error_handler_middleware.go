package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/shared"
)

var mapErrorCode = map[shared.CustomErrorCode]int{
	shared.BadRequest:     http.StatusBadRequest,
	shared.Forbidden:      http.StatusForbidden,
	shared.Unauthorized:   http.StatusUnauthorized,
	shared.InternalServer: http.StatusInternalServerError,
	shared.NotFound:       http.StatusNotFound,
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err != nil {
			switch e := err.Err.(type) {
			case *shared.CustomError:
				c.AbortWithStatusJSON(mapErrorCode[e.Code], e.CreateHTTPErrorMessage())
			case validator.ValidationErrors:
				errs := shared.ValidationErrResponse(e)
				c.AbortWithStatusJSON(http.StatusBadRequest, dto.JSONResponse{
					Message: errs,
				})
			default:
				if errors.Is(err, context.DeadlineExceeded) {
					c.AbortWithStatus(http.StatusRequestTimeout)
				} else {
					c.AbortWithStatusJSON(http.StatusInternalServerError, dto.JSONResponse{
						Message: "internal server error",
					})
				}
			}

			c.Abort()
		}
	}
}
