package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// clientLimiter controla tokens por cliente IP
type clientLimiter struct {
	tokens     int
	lastRefill time.Time
}

var (
	limiterMu    sync.Mutex
	clients      = map[string]*clientLimiter{}
	maxTokens    = 60 // requests
	refillPeriod = time.Minute
)

// RateLimitMiddleware aplica limitação simples por IP (maxTokens por refillPeriod)
func RateLimitMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			// remover porta se presente
			if strings.Contains(ip, ":") {
				parts := strings.Split(ip, ":")
				ip = strings.Join(parts[:len(parts)-1], ":")
			}

			limiterMu.Lock()
			cl, ok := clients[ip]
			if !ok {
				cl = &clientLimiter{tokens: maxTokens, lastRefill: time.Now()}
				clients[ip] = cl
			}
			// refill
			now := time.Now()
			elapsed := now.Sub(cl.lastRefill)
			if elapsed >= refillPeriod {
				cl.tokens = maxTokens
				cl.lastRefill = now
			}

			if cl.tokens <= 0 {
				limiterMu.Unlock()
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "rate limit exceeded",
				})
			}
			cl.tokens--
			limiterMu.Unlock()
			return next(c)
		}
	}
}
