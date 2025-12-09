package middleware

import (
	"net/http"

	"api-go-arquitetura/internal/logger"
)

// RecoveryMiddleware captura panics e retorna 500
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.WithFields(map[string]interface{}{
					"path":   r.URL.Path,
					"method": r.Method,
					"panic":  rec,
				}).Error("Panic recovered")
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
