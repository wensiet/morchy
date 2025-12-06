package infrastructure

import (
	ginzapcontrib "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/wernsiet/morchy/docs"
	"go.uber.org/zap"
)

// @title			Morchy
// @version		1.0
// @description	API for a distributed workloads management system
// @host			localhost:8080
// @BasePath		/
func NewRouter(logger *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(ginzapcontrib.RecoveryWithZap(logger, true))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
