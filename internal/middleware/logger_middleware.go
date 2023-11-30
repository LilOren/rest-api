package middleware

import (
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/lil-oren/rest/internal/dependency"
)

func Logger(log dependency.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		log.Info(map[string]interface{}{
			"id":      requestid.Get(c),
			"method":  reqMethod,
			"latency": latencyTime,
			"uri":     reqUri,
			"status":  statusCode,
			"ip":      clientIP,
		})

		c.Next()
	}
}
