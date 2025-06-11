package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto"
)

var HelloCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "hello_requests_total",
		Help: "Total number of /hello requests",
	},
)


