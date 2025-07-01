package api

import (
	"io"
	"math" // Import the math package
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"go.uber.org/zap"
)

// processBuckets is a helper function to make histogram buckets JSON-safe.
// It converts +Inf upper bounds to a string representation.
func processBuckets(buckets []*dto.Bucket) []gin.H {
	processed := make([]gin.H, len(buckets))
	for i, b := range buckets {
		upperBound := b.GetUpperBound()
		var boundValue interface{}

		// Check for positive infinity and convert it to a string for JSON compatibility.
		if math.IsInf(upperBound, 1) {
			boundValue = "+Inf"
		} else {
			boundValue = upperBound
		}

		processed[i] = gin.H{
			"cumulative_count": b.GetCumulativeCount(),
			"upper_bound":      boundValue,
		}
	}
	return processed
}

// Metrics godoc
// @Summary Get a specific metric or component of a metric.
// @Description Get a specific metric by name from the metrics endpoint. For simple counters and gauges, it returns the value. For complex types like Histograms and Summaries, you can specify a component (e.g., my_metric_sum, my_metric_count) or get a composite object by requesting the base name.
// @Tags metrics
// @Accept  json
// @Produce  json
// @Param name query string true "The full name of the metric to fetch (e.g., 'go_goroutines', 'http_request_duration_seconds_sum', 'http_request_duration_seconds_bucket')."
// @Success 200 {object} object "An object or an array of objects containing the metric data. For labeled metrics, an array is returned. For a single unlabeled metric, a single object is returned."
// @Failure 400 {object} gin.H "Error: name query parameter is required"
// @Failure 404 {object} gin.H "Error: metric not found"
// @Failure 500 {object} gin.H "Error: internal server error"
// @Router /api/v1/metrics [get]
func GetMetrics(c *gin.Context) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	metricsName := c.Query("name")
	if metricsName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'name' query parameter is required"})
		return
	}

	resp, err := http.Get("http://localhost:4000/metrics")
	if err != nil {
		logger.Error("failed to get metrics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get metrics"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("failed to read metrics response body", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read metrics response body"})
		return
	}

	parser := expfmt.TextParser{}
	metricFamilies, err := parser.TextToMetricFamilies(strings.NewReader(string(body)))
	if err != nil {
		logger.Error("failed to parse metrics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse metrics"})
		return
	}

	baseName, suffix := parseMetricName(metricsName)

	mf, ok := metricFamilies[baseName]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "metric not found"})
		return
	}

	results := []gin.H{}
	metricType := mf.GetType()

	for _, m := range mf.GetMetric() {
		labels := make(gin.H)
		for _, l := range m.GetLabel() {
			labels[l.GetName()] = l.GetValue()
		}

		var value interface{}

		switch metricType {
		case dto.MetricType_COUNTER:
			if suffix == "" {
				value = m.GetCounter().GetValue()
			}
		case dto.MetricType_GAUGE:
			if suffix == "" {
				value = m.GetGauge().GetValue()
			}
		case dto.MetricType_HISTOGRAM:
			h := m.GetHistogram()
			switch suffix {
			case "_sum":
				value = h.GetSampleSum()
			case "_count":
				value = h.GetSampleCount()
			case "_bucket":
				// Use the helper function to process buckets before assigning
				value = processBuckets(h.GetBucket())
			case "":
				// Use the helper function for the composite object as well
				value = gin.H{
					"count":   h.GetSampleCount(),
					"sum":     h.GetSampleSum(),
					"buckets": processBuckets(h.GetBucket()),
				}
			}
		case dto.MetricType_SUMMARY:
			s := m.GetSummary()
			switch suffix {
			case "_sum":
				value = s.GetSampleSum()
			case "_count":
				value = s.GetSampleCount()
			case "":
				value = gin.H{
					"count":     s.GetSampleCount(),
					"sum":       s.GetSampleSum(),
					"quantiles": s.GetQuantile(),
				}
			}
		}

		if value != nil {
			metricData := gin.H{"value": value}
			if len(labels) > 0 {
				metricData["labels"] = labels
			}
			results = append(results, metricData)
		}
	}

	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "metric component not found or suffix is not applicable for this metric type"})
		return
	}
	
	if len(results) == 1 {
		labels, ok := results[0]["labels"].(gin.H)
		if !ok || len(labels) == 0 {
			c.JSON(http.StatusOK, results[0])
			return
		}
	}
	
	c.JSON(http.StatusOK, results)
}

func parseMetricName(metricName string) (string, string) {
	suffixes := []string{"_bucket", "_sum", "_count"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(metricName, suffix) {
			return strings.TrimSuffix(metricName, suffix), suffix
		}
	}
	return metricName, ""
}