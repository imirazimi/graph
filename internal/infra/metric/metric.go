package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "status"},
	)

	RequestLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_latency_histogram",
			Help:    "Request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	TasksCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "tasks_count",
			Help: "Current number of tasks",
		},
	)
)

func init() {
	prometheus.MustRegister(
		RequestsTotal,
		RequestLatencyHistogram,
		TasksCount,
	)
}
