package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wernsiet/morchy/pkg/infrastructure/mtls"
	"go.uber.org/zap"
)

const (
	HeaderNodeID = "X-Client-Node-ID"
)

func MTLSAuthMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.TLS == nil {
			logger.Warn("request without TLS connection")
			c.AbortWithStatusJSON(400, gin.H{"error": "TLS connection required"})
			return
		}

		if len(c.Request.TLS.PeerCertificates) == 0 {
			logger.Warn("request without client certificate")
			c.AbortWithStatusJSON(401, gin.H{"error": "client certificate required"})
			return
		}

		clientCert := c.Request.TLS.PeerCertificates[0]

		nodeID, err := mtls.ExtractNodeIDFromCertificate(clientCert)
		if err != nil {
			logger.Error("failed to extract node ID from certificate", zap.Error(err))
			c.AbortWithStatusJSON(401, gin.H{"error": fmt.Sprintf("invalid certificate: %v", err)})
			return
		}

		c.Set("node_id", nodeID)
		c.Set("client_cert", clientCert)
		c.Request.Header.Set(HeaderNodeID, nodeID)

		c.Next()
	}
}

func GetNodeIDFromContext(c *gin.Context) (string, error) {
	nodeID, exists := c.Get("node_id")
	if !exists {
		return "", fmt.Errorf("node_id not found in context")
	}

	nodeIDStr, ok := nodeID.(string)
	if !ok {
		return "", fmt.Errorf("node_id is not a string")
	}

	return nodeIDStr, nil
}
