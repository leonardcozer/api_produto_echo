package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"api-go-arquitetura/internal/errors"
)

// JSONResponse envia uma resposta JSON com status code
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Se houver erro ao codificar JSON, enviar erro genérico
		http.Error(w, "Erro ao processar resposta", http.StatusInternalServerError)
	}
}

// SuccessResponse envia uma resposta de sucesso
func SuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	JSONResponse(w, status, data)
}

// ErrorResponse envia uma resposta de erro padronizada
func ErrorResponse(w http.ResponseWriter, err error) {
	// Verificar se é um erro customizado da API
	if apiErr := errors.AsAPIError(err); apiErr != nil {
		JSONResponse(w, apiErr.Status, apiErr)
		return
	}

	// Se não for, tratar como erro interno
	apiErr := errors.ErrInternalServer.WithDetails(err.Error())
	JSONResponse(w, apiErr.Status, apiErr)
}

// ValidationErrorResponse envia uma resposta de erro de validação
func ValidationErrorResponse(w http.ResponseWriter, validationErrors []string) {
	apiErr := errors.ErrValidation.WithDetailsf("Erros de validação: %v", validationErrors)
	JSONResponse(w, apiErr.Status, apiErr)
}

// NotFoundResponse envia uma resposta de recurso não encontrado
func NotFoundResponse(w http.ResponseWriter, resource string) {
	apiErr := errors.ErrNotFound.WithDetailsf("%s não encontrado", resource)
	JSONResponse(w, apiErr.Status, apiErr)
}

// BadRequestResponse envia uma resposta de requisição inválida
func BadRequestResponse(w http.ResponseWriter, message string) {
	apiErr := errors.ErrInvalidInput.WithDetails(message)
	JSONResponse(w, apiErr.Status, apiErr)
}

// DecodeJSON decodifica um JSON do body da requisição
func DecodeJSON(body io.Reader, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}

