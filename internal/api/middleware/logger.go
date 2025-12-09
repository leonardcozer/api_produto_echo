package middleware

import (
	"net/http"
	"time"

	"api-go-arquitetura/internal/logger"
)

// responseWriter é um wrapper para http.ResponseWriter que captura o status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware registra solicitações com método, path, remote addr, status e duração
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)
		
		// Obter request ID do contexto
		requestID := GetRequestID(r)
		
		// Log da requisição recebida
		logFields := map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		}
		if requestID != "" {
			logFields["request_id"] = requestID
		}
		logger.WithFields(logFields).Info("Request received")
		
		next.ServeHTTP(rw, r)
		
		dur := time.Since(start)
		
		// Log da resposta
		responseFields := map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status_code": rw.statusCode,
			"duration_ms": dur.Milliseconds(),
			"duration":    dur.String(),
		}
		if requestID != "" {
			responseFields["request_id"] = requestID
		}
		logger.WithFields(responseFields).Info("Request completed")
	})
}
