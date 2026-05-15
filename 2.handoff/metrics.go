package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "handoff_http_requests_total",
			Help: "	Total requests",
		},
		[]string{"method", "path", "status_code"},
	)

	httpDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "handoff_http_request_duration_seconds",
			Help: "Request latency distribution",
			// Buckets: prometheus.DefBuckets,
			Buckets: []float64{.03, .1},
		},
		[]string{"method", "path"},
	)

	incidentTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "handoff_incidents_total",
			Help: "Current number of incidents",
		},
		[]string{"status"},
	)

	totalEntries = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "handoff_entries_total",
			Help: "Total timeline entries created",
		},
	)

	dbQueryDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "handoff_db_query_duration_seconds",
			Help:    "Database query latency",
			Buckets: []float64{.03, .1},
		},
		[]string{"operation"},
	)

	wsConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "handoff_websocket_connections",
			Help: "Current number of active WebSocket connections",
		},
	)
)

func NewMetrics(reg *prometheus.Registry) {
	reg.MustRegister(httpRequestsTotal)
	reg.MustRegister(httpDurationSeconds)
	reg.MustRegister(incidentTotal)
	reg.MustRegister(totalEntries)
	reg.MustRegister(dbQueryDurationSeconds)
	reg.MustRegister(wsConnections)
}
