package middleware

import (
	"crypto/x509"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MTLSConfig struct {
	DevMode bool
}

func NewMTLSMiddleware(config MTLSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.DevMode {
			c.Next()
			return
		}

		if c.Request.TLS == nil || len(c.Request.TLS.PeerCertificates) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Client certificate required"})
			c.Abort()
			return
		}

		cert := c.Request.TLS.PeerCertificates[0]
		if cert == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid client certificate"})
			c.Abort()
			return
		}

		if len(cert.Subject.CommonName) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Certificate must have Common Name"})
			c.Abort()
			return
		}

		c.Set("node_id", cert.Subject.CommonName)
		c.Next()
	}
}

func ExtractNodeIDFromCert(c *gin.Context) string {
	if nodeID, exists := c.Get("node_id"); exists {
		if id, ok := nodeID.(string); ok {
			return id
		}
	}
	return ""
}

func VerifyCertChain(cert *x509.Certificate, roots *x509.CertPool) bool {
	if cert == nil || roots == nil {
		return false
	}

	opts := x509.VerifyOptions{
		Roots: roots,
	}

	_, err := cert.Verify(opts)
	return err == nil
}
