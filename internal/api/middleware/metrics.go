package middleware

import (
	"time"

	"github.com/labstack/echo/v4"

	"api-go-arquitetura/internal/metrics"
)

// MetricsMiddleware registra métricas Prometheus
func MetricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Incrementar conexões ativas
			metrics.ActiveConnections.Inc()
			defer metrics.ActiveConnections.Dec()

			err := next(c)

			duration := time.Since(start)
			statusCode := c.Response().Status

			// Registrar métricas
			metrics.RecordHTTPRequest(c.Request().Method, c.Request().URL.Path, statusCode, duration)

			return err
		}
	}
}
