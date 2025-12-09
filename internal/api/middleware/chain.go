package middleware

import (
	"github.com/labstack/echo/v4"
)

// ApplyMiddlewares aplica a cadeia de middlewares ao Echo
// Ordem: RequestID -> Metrics -> Logging -> Recovery -> CORS -> RateLimit
func ApplyMiddlewares(e *echo.Echo) {
	// Echo já tem middlewares built-in, então vamos usar a ordem correta
	e.Use(RequestIDMiddleware())
	e.Use(MetricsMiddleware())
	e.Use(LoggingMiddleware())
	e.Use(RecoveryMiddleware())
	e.Use(CORSMiddleware())
	e.Use(RateLimitMiddleware())
}
