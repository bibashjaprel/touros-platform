package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	permitsIssuedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "permits_issued_total",
			Help: "Total number of permits issued",
		},
	)

	checkInsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "check_ins_total",
			Help: "Total number of safety check-ins",
		},
	)

	sosIncidentsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "sos_incidents_total",
			Help: "Total number of SOS incidents",
		},
	)
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Observe(duration)
	}
}

func IncrementPermitsIssued() {
	permitsIssuedTotal.Inc()
}

func IncrementCheckIns() {
	checkInsTotal.Inc()
}

func IncrementSOSIncidents() {
	sosIncidentsTotal.Inc()
}

