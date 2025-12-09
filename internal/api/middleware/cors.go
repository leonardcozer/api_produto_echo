package middleware

import (
	"net/http"
	"strings"

	"api-go-arquitetura/internal/config"
)

var corsConfig *config.Config

// SetCORSConfig configura o middleware CORS
func SetCORSConfig(cfg *config.Config) {
	corsConfig = cfg
}

// CORSMiddleware adiciona cabeçalhos CORS e responde OPTIONS
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Se não há configuração, usar padrão permissivo
		if corsConfig == nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		} else {
			// Configurar origem
			origin := r.Header.Get("Origin")
			allowedOrigin := getAllowedOrigin(origin)
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)

			// Configurar métodos
			if len(corsConfig.CORSAllowedMethods) > 0 {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.CORSAllowedMethods, ", "))
			} else {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			}

			// Configurar headers
			if len(corsConfig.CORSAllowedHeaders) > 0 {
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.CORSAllowedHeaders, ", "))
			} else {
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			// Configurar credenciais
			if corsConfig.CORSCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		// Responder a requisições OPTIONS
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getAllowedOrigin retorna a origem permitida baseada na configuração
func getAllowedOrigin(requestOrigin string) string {
	if corsConfig == nil || len(corsConfig.CORSAllowedOrigins) == 0 {
		return "*"
	}

	// Se "*" está na lista, permitir todas
	for _, origin := range corsConfig.CORSAllowedOrigins {
		if origin == "*" {
			return "*"
		}
		if origin == requestOrigin {
			return origin
		}
	}

	// Se não encontrou, retornar a primeira origem permitida (ou "*" se vazio)
	if len(corsConfig.CORSAllowedOrigins) > 0 {
		return corsConfig.CORSAllowedOrigins[0]
	}
	return "*"
}
