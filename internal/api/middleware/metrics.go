package middleware

import (
	"net/http"
	"time"

	"api-go-arquitetura/internal/metrics"
)

// MetricsMiddleware registra métricas Prometheus
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)

		// Incrementar conexões ativas
		metrics.ActiveConnections.Inc()
		defer metrics.ActiveConnections.Dec()

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Registrar métricas
		metrics.RecordHTTPRequest(r.Method, r.URL.Path, rw.statusCode, duration)
	})
}
