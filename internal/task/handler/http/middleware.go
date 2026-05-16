package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imirazimi/graph/internal/infra/metric"
)

func MetricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start).Seconds()

		method := c.Request.Method
		route := c.FullPath()
		status := strconv.Itoa(c.Writer.Status())

		metric.RequestsTotal.WithLabelValues(
			method,
			route,
			status,
		).Inc()

		metric.RequestLatencyHistogram.WithLabelValues(
			method,
			route,
		).Observe(latency)
	}
}