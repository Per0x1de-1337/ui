package telemetry

import (
	"os/exec"

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

var (
	// Counter metrics for operations
	BindingPolicyOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kubestellar_binding_policy_operations_total",
			Help: "Total number of binding policy operations",
		},
		[]string{"operation", "status"},
	)

	// Histogram for operation latency
	BindingPolicyOperationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kubestellar_binding_policy_operation_duration_seconds",
			Help:    "Duration of binding policy operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// Cache hit/miss ratios
	BindingPolicyCacheHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kubestellar_binding_policy_cache_hits_total",
			Help: "Total cache hits for binding policies",
		},
		[]string{"cache_type"},
	)

	BindingPolicyCacheMisses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kubestellar_binding_policy_cache_misses_total",
			Help: "Total cache misses for binding policies",
		},
		[]string{"cache_type"},
	)

	BindingPolicyWatchEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kubestellar_binding_policy_watch_events_total",
			Help: "Total watch events processed",
		},
		[]string{"event_type", "status"},
	)

	// Reconciliation time tracking
	BindingPolicyReconciliationDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "kubestellar_binding_policy_reconciliation_duration_seconds",
			Help:    "Time taken for binding policy reconciliation",
			Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0, 10.0, 30.0},
		},
	)

	ClusterOnboardingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cluster_onboarding_duration_seconds",
			Help:    "Duration of cluster onboarding process",
			Buckets: []float64{30, 60, 120, 300, 600, 900, 1800}, // 30s to 30min
		},
		[]string{"cluster_name", "status"},
	)
	KubectlOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kubectl_operations_total",
			Help: "Total number of kubectl operations executed",
		},
		[]string{"command", "context", "status"},
	)

	GithubDeploymentsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "github_deployments_total",
			Help: "Total number of GitHub deployments created",
		},
		[]string{"type", "status"},
	)
	WebsocketConnectionsActive = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "websocket_connections_active",
			Help: "Number of active WebSocket connections",
		},
		[]string{"endpoint", "type"},
	)

	WebsocketConnectionsFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_connections_failed_total",
			Help: "Total number of failed WebSocket connections",
		},
		[]string{"endpoint", "error_type"},
	)

	WebsocketConnectionUpgradedSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_connection_upgraded_success_total",
			Help: "Total number of successful WebSocket connection upgrades",
		},
		[]string{"endpoint", "type"},
	)
	WebsocketConnectionUpgradedFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_connection_upgraded_failed_total",
			Help: "Total number of failed WebSocket connection upgrades",
		},
		[]string{"endpoint", "error_type"})
)

func InstrumentKubectlCommand(cmd *exec.Cmd, operation string, context string) {
	// start := time.Now()
	// err := cmd.Run()

	// status := "success"
	// if err != nil {
	// 	status = "failed"
	// }

	KubectlOperationsTotal.WithLabelValues(operation, context).Inc()
}
