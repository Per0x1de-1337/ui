package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto"
)

var TotalHTTPRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "path", "status_code"},
)

var HTTPRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path"},
)

// var HTTPResponseSize = prometheus.NewHistogramVec(
// 	prometheus.HistogramOpts{
// 		Name:    "http_response_size_bytes",
// 		Help:    "Size of HTTP responses in bytes",
// 		Buckets: prometheus.ExponentialBuckets(100, 2, 10), // 100 bytes to 102400 bytes
// 	},
// 	[]string{"method", "path", "status_code"},
// )

// var HTTPRequestSize = prometheus.NewHistogramVec(
// 	prometheus.HistogramOpts{
// 		Name:    "http_request_size_bytes",
// 		Help:    "Size of HTTP requests in bytes",
// 		Buckets: prometheus.ExponentialBuckets(100, 2, 10), // 100 bytes to 102400 bytes
// 	},
// 	[]string{"method", "path"},
// )

var HTTPErrorCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_error_requests_total",
		Help: "Total number of HTTP error requests",
	},
	[]string{"method", "path", "status_code"},
)
