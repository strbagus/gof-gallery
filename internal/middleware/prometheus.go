package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests.",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "status"})
)

func Prometheus(c fiber.Ctx) error {
	if c.Path() == "/metrics" {
		return c.Next()
	}

	start := time.Now()

	err := c.Next()

	path := c.FullPath()
	if path == "" {
		path = c.Path()
	}

	status := c.Response().StatusCode()
	duration := time.Since(start).Seconds()

	statusStr := strconv.Itoa(status)
	method := c.Method()

	httpRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
	httpRequestDuration.WithLabelValues(method, path, statusStr).Observe(duration)

	return err
}
