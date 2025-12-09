package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
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
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar se já existe um request ID no header da requisição
		requestID := r.Header.Get(RequestIDHeader)
		
		// Se não existir, gerar um novo UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// Adicionar o request ID no header de resposta
		w.Header().Set(RequestIDHeader, requestID)
		
		// Adicionar o request ID no contexto da requisição
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)
		
		// Continuar com o próximo handler
		next.ServeHTTP(w, r)
	})
}

// GetRequestID extrai o request ID do contexto da requisição
func GetRequestID(r *http.Request) string {
	if requestID, ok := r.Context().Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

