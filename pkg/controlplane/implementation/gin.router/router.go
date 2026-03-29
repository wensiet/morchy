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

func (rh *RouterHandler) SetRoutes(r *gin.Engine, seedTokenCfg middlewareauth.SeedTokenConfig, mtlsCfg middlewareauth.MTLSConfig, dualAuthCfg middlewareauth.DualAuthConfig) {
	apiV1 := r.Group("/api/v1")
	rh.setWorkloadRoutes(apiV1, seedTokenCfg, mtlsCfg, dualAuthCfg)
}

func (rh *RouterHandler) setWorkloadRoutes(apiV1 *gin.RouterGroup, seedTokenCfg middlewareauth.SeedTokenConfig, mtlsCfg middlewareauth.MTLSConfig, dualAuthCfg middlewareauth.DualAuthConfig) {
	seedTokenAuth := middlewareauth.NewSeedTokenMiddleware(seedTokenCfg)
	mtlsAuth := middlewareauth.NewMTLSMiddleware(mtlsCfg)
	dualAuth := middlewareauth.NewDualAuthMiddleware(dualAuthCfg)

	dualAuthGroup := apiV1.Group("")
	dualAuthGroup.Use(dualAuth)
	{
		dualAuthGroup.GET("/workloads", rh.listWorkloads)
	}

	adminGroup := apiV1.Group("")
	adminGroup.Use(seedTokenAuth)
	{
		adminGroup.GET("/workloads/:workload_id", rh.getWorkload)
		adminGroup.POST("/workloads", rh.createWorkload)
		adminGroup.DELETE("/workloads/:workload_id", rh.deleteWorkload)
		adminGroup.GET("/edges", rh.listEdges)
	}

	agentGroup := apiV1.Group("")
	agentGroup.Use(mtlsAuth)
	{
		agentGroup.GET("/workloads/:workload_id/lease", rh.getLease)
		agentGroup.PUT("/workloads/:workload_id/lease", rh.putLease)
		agentGroup.DELETE("/workloads/:workload_id/lease", rh.deleteLease)
		agentGroup.POST("/events", rh.pushEvent)
	}
}
