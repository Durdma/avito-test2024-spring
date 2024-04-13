package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"method"},
	)

	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)
)

func Init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(responseStatus)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		method := c.Request.Method
		elapsed := time.Since(start).Seconds()
		requestDuration.WithLabelValues(method).Observe(elapsed)

		statusCode := c.Writer.Status()
		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
	}
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
