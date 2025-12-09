package middleware

import (
	"time"

	"github.com/labstack/echo/v4"

	"api-go-arquitetura/internal/logger"
)

// LoggingMiddleware registra solicitações com método, path, remote addr, status e duração
func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			// Obter request ID do contexto
			requestID := GetRequestID(c)
			
			// Log da requisição recebida
			logFields := map[string]interface{}{
				"method":      c.Request().Method,
				"path":        c.Request().URL.Path,
				"remote_addr": c.Request().RemoteAddr,
				"user_agent":  c.Request().UserAgent(),
			}
			if requestID != "" {
				logFields["request_id"] = requestID
			}
			logger.WithFields(logFields).Info("Request received")
			
			err := next(c)
			
			dur := time.Since(start)
			statusCode := c.Response().Status
			
			// Log da resposta
			responseFields := map[string]interface{}{
				"method":      c.Request().Method,
				"path":        c.Request().URL.Path,
				"status_code": statusCode,
				"duration_ms": dur.Milliseconds(),
				"duration":    dur.String(),
			}
			if requestID != "" {
				responseFields["request_id"] = requestID
			}
			logger.WithFields(responseFields).Info("Request completed")
			
			return err
		}
	}
}
