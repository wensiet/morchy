package ginrouter

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (rh *RouterHandler) proxyHandler(c *gin.Context) {
	path := c.Request.URL.Path

	pathList := strings.Split(path, "/")
	if len(pathList) < 2 {
		rh.logger.Warn("path is empty", zap.String("path", path))
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	edge := rh.ucHandler.FindEdge(c.Request.Context(), "/"+pathList[1])
	rh.logger.Info("request path is: ", zap.String("request_path", path))
	if edge == nil {
		rh.logger.Debug("no edge found for path", zap.String("path", path))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	backendURL := "http://" + edge.UpstreamAddress
	if backendURL == "" {
		rh.logger.Warn("edge has empty backend URL", zap.String("path", path))
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	target, err := url.Parse(backendURL)
	if err != nil {
		rh.logger.Warn("invalid backend URL",
			zap.String("url", backendURL),
			zap.Error(err))
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	prefix := edge.ProxyPath
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	c.Request.URL.Path = strings.TrimPrefix(path, prefix)

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(c.Writer, c.Request)
}
