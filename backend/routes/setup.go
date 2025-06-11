package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kubestellar/ui/plugin/plugins"
	"github.com/kubestellar/ui/telemetry"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	fmt.Println("Registering Prometheus metrics..(((((((((((((((((((((((()))))))))))))))))))))))).")
	prometheus.MustRegister(telemetry.TotalHTTPRequests)
	prometheus.MustRegister(telemetry.HTTPRequestDuration)
	prometheus.MustRegister(telemetry.HTTPErrorCounter)
}
func SetupRoutes(router *gin.Engine) {
	// Initialize all route groups
	setupClusterRoutes(router)
	setupDeploymentRoutes(router)
	setupNamespaceRoutes(router)
	setupBindingPolicyRoutes(router)
	setupResourceRoutes(router)
	getWecsResources(router)
	setupInstallerRoutes(router)
	setupWdsCookiesRoute(router)
	setupGitopsRoutes(router)
	setupHelmRoutes(router)
	setupGitHubRoutes(router)
	setupDeploymentHistoryRoutes(router)
	plugins.Pm.SetupPluginsRoutes(router)

	setupAuthRoutes(router)
	setupArtifactHubRoutes(router)
	setupMetricsRoutes(router)
}
