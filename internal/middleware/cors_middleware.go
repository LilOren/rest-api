package middleware

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lil-oren/rest/internal/dependency"
)

func CORS(config dependency.Config) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowMethods:     []string{"POST", "PUT", "GET", "PATCH", "DELETE"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
	}

	if config.App.OriginDomain == "localhost" {
		origin := "http://localhost:3000"
		corsConfig.AllowOrigins = []string{origin}
	} else {
		origin := fmt.Sprintf("https://%s", config.App.OriginDomain)
		corsConfig.AllowOrigins = []string{origin}
	}

	return cors.New(corsConfig)
}
