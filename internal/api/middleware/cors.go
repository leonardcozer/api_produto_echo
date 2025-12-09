package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"api-go-arquitetura/internal/config"
)

var corsConfig *config.Config

// SetCORSConfig configura o middleware CORS
func SetCORSConfig(cfg *config.Config) {
	corsConfig = cfg
}

// CORSMiddleware adiciona cabeçalhos CORS e responde OPTIONS
func CORSMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Se não há configuração, usar padrão permissivo
			if corsConfig == nil {
				c.Response().Header().Set("Access-Control-Allow-Origin", "*")
				c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			} else {
				// Configurar origem
				origin := c.Request().Header.Get("Origin")
				allowedOrigin := getAllowedOrigin(origin)
				c.Response().Header().Set("Access-Control-Allow-Origin", allowedOrigin)

				// Configurar métodos
				if len(corsConfig.CORSAllowedMethods) > 0 {
					c.Response().Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.CORSAllowedMethods, ", "))
				} else {
					c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				}

				// Configurar headers
				if len(corsConfig.CORSAllowedHeaders) > 0 {
					c.Response().Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.CORSAllowedHeaders, ", "))
				} else {
					c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				}

				// Configurar credenciais
				if corsConfig.CORSCredentials {
					c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			// Responder a requisições OPTIONS
			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}

			return next(c)
		}
	}
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
