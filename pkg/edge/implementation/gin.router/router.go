package ginrouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wernsiet/morchy/pkg/edge/usecase"
	"go.uber.org/zap"
)

type RouterHandler struct {
	logger    *zap.Logger
	ucHandler usecase.Handler
}

func NewRouterHandler(logger *zap.Logger, ucHandler usecase.Handler) RouterHandler {
	return RouterHandler{
		logger:    logger,
		ucHandler: ucHandler,
	}
}

func (rh *RouterHandler) SetRoutes(r *gin.Engine) {
	r.NoRoute(rh.proxyHandler)
}
