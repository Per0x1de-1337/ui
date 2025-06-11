package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto"
)

var TotalWebSocketConnections = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "websocket_connections_total",
		Help: "Total number of WebSocket connections",
	},
)

var TotalActiveWebSocketConnections = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "active_websocket_connections",
		Help: "Current number of active WebSocket connections",
	},
)

var TotalWebSocketsSentMessages = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "websocket_sent_messages_total",
		Help: "Total number of messages sent over WebSocket connections",
	},
	[]string{"path"},
)
var TotalWebSocketsReceivedMessages = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "websocket_received_messages_total",
		Help: "Total number of messages received over WebSocket connections",
	},
	[]string{"path"},
)

var TotalDurationWebSocketMessages = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "websocket_message_duration_seconds",
		Help:    "Duration of WebSocket messages in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"path"},
)


