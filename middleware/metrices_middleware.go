package middleware

import (
	"net/http"
	"time"

	"github.com/m3db/prometheus_client_golang/prometheus"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of http requests",
		},
		[]string{"path", "method"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_requests_duration_seconds",
			Help: "Duration of http requests in seconds",
		},
		[]string{"path", "method"},
	)

	statusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_total",
			Help: "Total Number of http response by status code",
		},
		[]string{"path", "method", "status_code"},
	)
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func MetricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start time of the request
		start := time.Now()

		ww := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(ww, r)

		// Duration of the request
		duration := time.Since(start).Seconds()

		requestCounter.WithLabelValues(r.URL.Path, r.Method).Inc()

		requestDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration)

		statusCounter.WithLabelValues(r.URL.Path, r.Method, http.StatusText(ww.statusCode))
	})
}

func init() {
	prometheus.MustRegister(requestCounter, requestDuration, statusCounter)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
