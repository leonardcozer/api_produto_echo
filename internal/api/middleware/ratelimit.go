package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
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
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprint(w, `{"error":"rate limit exceeded"}`)
			return
		}
		cl.tokens--
		limiterMu.Unlock()
		next.ServeHTTP(w, r)
	})
}
