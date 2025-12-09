package middleware

import (
	"net/http"
)

// ApplyMiddlewares aplica a cadeia de middlewares ao handler fornecido
// Ordem: RequestID -> Metrics -> Logging -> Recovery -> CORS -> RateLimit
func ApplyMiddlewares(h http.Handler) http.Handler {
	return RateLimitMiddleware(CORSMiddleware(RecoveryMiddleware(LoggingMiddleware(MetricsMiddleware(RequestIDMiddleware(h))))))
}
