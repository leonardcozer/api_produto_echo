package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ContextKey é o tipo usado para chaves do contexto
type contextKey string

const (
	// RequestIDKey é a chave usada para armazenar o request ID no contexto
	RequestIDKey contextKey = "request_id"
	// RequestIDHeader é o nome do header HTTP usado para request ID
	RequestIDHeader = "X-Request-ID"
)

// RequestIDMiddleware gera um ID único para cada requisição e o adiciona:
// - No header de resposta (X-Request-ID)
// - No contexto da requisição
// - Nos logs (através do contexto)
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Verificar se já existe um request ID no header da requisição
			requestID := c.Request().Header.Get(RequestIDHeader)
			
			// Se não existir, gerar um novo UUID
			if requestID == "" {
				requestID = uuid.New().String()
			}
			
			// Adicionar o request ID no header de resposta
			c.Response().Header().Set(RequestIDHeader, requestID)
			
			// Adicionar o request ID no contexto da requisição
			ctx := context.WithValue(c.Request().Context(), RequestIDKey, requestID)
			c.SetRequest(c.Request().WithContext(ctx))
			
			// Continuar com o próximo handler
			return next(c)
		}
	}
}

// GetRequestID extrai o request ID do contexto da requisição
func GetRequestID(c echo.Context) string {
	if requestID, ok := c.Request().Context().Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

