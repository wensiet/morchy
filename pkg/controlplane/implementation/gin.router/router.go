package ginrouter

import (
	"github.com/gin-gonic/gin"
	middlewareauth "github.com/wernsiet/morchy/pkg/controlplane/infrastructure/middleware"
	"github.com/wernsiet/morchy/pkg/controlplane/usecase"
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

func (rh *RouterHandler) SetRoutes(r *gin.Engine, dualAuthCfg middlewareauth.DualAuthConfig) {
	dualAuth := middlewareauth.NewDualAuthMiddleware(dualAuthCfg)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(dualAuth)
	{
		apiV1.GET("/workloads", rh.listWorkloads)
		apiV1.GET("/workloads/:workload_id", rh.getWorkload)
		apiV1.POST("/workloads", rh.createWorkload)
		apiV1.DELETE("/workloads/:workload_id", rh.deleteWorkload)
		apiV1.GET("/workloads/:workload_id/lease", rh.getLease)
		apiV1.PUT("/workloads/:workload_id/lease", rh.putLease)
		apiV1.DELETE("/workloads/:workload_id/lease", rh.deleteLease)
		apiV1.POST("/events", rh.pushEvent)
		apiV1.GET("/edges", rh.listEdges)
	}
}
