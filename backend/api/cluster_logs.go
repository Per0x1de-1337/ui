package api

import (
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/kubestellar/ui/telemetry"
)

// OnboardingLogsHandler returns all logs for a specific cluster's onboarding process
func OnboardingLogsHandler(c *gin.Context) {
	clusterName := c.Param("cluster")
	startTime := time.Now()
	if clusterName == "" {
		telemetry.HTTPErrorCounter.WithLabelValues("GET", "/clusters/onboard/logs/:cluster", "400").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cluster name is required"})
		return
	}

	// Get all events for this cluster
	events := GetOnboardingEvents(clusterName)

	// Get current status
	mutex.RLock()
	status, exists := clusterStatuses[clusterName]
	mutex.RUnlock()

	if !exists {
		telemetry.HTTPErrorCounter.WithLabelValues("GET", "/clusters/onboard/logs/:cluster", "404").Inc()
		c.JSON(http.StatusNotFound, gin.H{"error": "No onboarding data found for cluster"})
		return
	}
	telemetry.HTTPRequestDuration.WithLabelValues("GET", "/clusters/onboard/logs/:cluster").Observe(time.Since(startTime).Seconds())
	telemetry.TotalHTTPRequests.WithLabelValues("GET", "/clusters/onboard/logs/:cluster", "200").Inc()
	c.JSON(http.StatusOK, gin.H{
		"clusterName": clusterName,
		"status":      status,
		"logs":        events,
		"count":       len(events),
	})
}
