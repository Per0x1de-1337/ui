package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kubestellar/ui/api"
)

func SetupMetricsRoutes(router *gin.Engine) {
	router.GET("/api/v1/metrics", api.GetMetrics)
}