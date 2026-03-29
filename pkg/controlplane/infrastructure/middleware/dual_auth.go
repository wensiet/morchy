package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type DualAuthConfig struct {
	SeedToken string
	DevMode   bool
}

func NewDualAuthMiddleware(config DualAuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.DevMode {
			c.Next()
			return
		}

		hasValidCert := c.Request.TLS != nil && len(c.Request.TLS.PeerCertificates) > 0 && c.Request.TLS.PeerCertificates[0] != nil && len(c.Request.TLS.PeerCertificates[0].Subject.CommonName) > 0
		hasValidToken := false

		token := c.GetHeader("X-Seed-Token")
		if token != "" && strings.EqualFold(token, config.SeedToken) {
			hasValidToken = true
		}

		if !hasValidCert && !hasValidToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required: either valid client certificate or seed token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
