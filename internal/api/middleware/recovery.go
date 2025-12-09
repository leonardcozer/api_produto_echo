package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"api-go-arquitetura/internal/logger"
)

// RecoveryMiddleware captura panics e retorna 500
func RecoveryMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if rec := recover(); rec != nil {
					logger.WithFields(map[string]interface{}{
						"path":   c.Request().URL.Path,
						"method": c.Request().Method,
						"panic":  rec,
					}).Error("Panic recovered")
					c.JSON(http.StatusInternalServerError, map[string]string{
						"error": "internal server error",
					})
				}
			}()
			return next(c)
		}
	}
}
