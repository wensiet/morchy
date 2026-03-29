package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SeedTokenConfig struct {
	SeedToken string
	DevMode   bool
}

func NewSeedTokenMiddleware(config SeedTokenConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.DevMode {
			c.Next()
			return
		}

		token := c.GetHeader("X-Seed-Token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing seed token"})
			c.Abort()
			return
		}

		if !strings.EqualFold(token, config.SeedToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid seed token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
