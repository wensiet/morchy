package ginrouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wernsiet/morchy/internal/usecase"
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
	apiV1 := r.Group("/api/v1")
	rh.setWorkloadRoutes(*apiV1)
}

func (rh *RouterHandler) setWorkloadRoutes(apiV1 gin.RouterGroup) {
	apiV1.GET("/workloads", rh.listWorkloads)
	apiV1.GET("/workloads/:workload_id", rh.getWorkload)
	apiV1.POST("/workloads", rh.createWorkload)
	apiV1.POST("/workloads/:workload_id/lease", rh.createLease)
	apiV1.PUT("/workloads/:workload_id/lease", rh.extendLease)
}
