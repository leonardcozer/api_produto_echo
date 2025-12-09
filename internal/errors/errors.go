package errors

import (
	"fmt"
	"net/http"
)

// APIError representa um erro padronizado da API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"-"` // Não serializa, usado apenas internamente
}

// Error implementa a interface error
func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// WithDetails adiciona detalhes ao erro
func (e *APIError) WithDetails(details string) *APIError {
	return &APIError{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
		Status:  e.Status,
	}
}

// WithDetailsf adiciona detalhes formatados ao erro
func (e *APIError) WithDetailsf(format string, args ...interface{}) *APIError {
	return e.WithDetails(fmt.Sprintf(format, args...))
}

// Erros pré-definidos da API
var (
	// Erros de validação (400)
	ErrInvalidInput = &APIError{
		Code:    "INVALID_INPUT",
		Message: "Dados de entrada inválidos",
		Status:  http.StatusBadRequest,
	}

	ErrInvalidID = &APIError{
		Code:    "INVALID_ID",
		Message: "ID inválido",
		Status:  http.StatusBadRequest,
	}

	// Erros de recurso não encontrado (404)
	ErrNotFound = &APIError{
		Code:    "NOT_FOUND",
		Message: "Recurso não encontrado",
		Status:  http.StatusNotFound,
	}

	ErrProdutoNotFound = &APIError{
		Code:    "PRODUTO_NOT_FOUND",
		Message: "Produto não encontrado",
		Status:  http.StatusNotFound,
	}

	// Erros de servidor (500)
	ErrInternalServer = &APIError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Erro interno do servidor",
		Status:  http.StatusInternalServerError,
	}

	ErrDatabase = &APIError{
		Code:    "DATABASE_ERROR",
		Message: "Erro ao acessar banco de dados",
		Status:  http.StatusInternalServerError,
	}

	// Erros de validação de negócio (422)
	ErrValidation = &APIError{
		Code:    "VALIDATION_ERROR",
		Message: "Erro de validação",
		Status:  http.StatusUnprocessableEntity,
	}

	ErrNomeObrigatorio = &APIError{
		Code:    "NOME_OBRIGATORIO",
		Message: "Nome do produto é obrigatório",
		Status:  http.StatusUnprocessableEntity,
	}

	ErrPrecoInvalido = &APIError{
		Code:    "PRECO_INVALIDO",
		Message: "Preço não pode ser negativo",
		Status:  http.StatusUnprocessableEntity,
	}
)

// IsAPIError verifica se um erro é do tipo APIError
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

// AsAPIError converte um erro para APIError, retornando nil se não for
func AsAPIError(err error) *APIError {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr
	}
	return nil
}

// WrapError envolve um erro genérico em um APIError
func WrapError(err error, apiErr *APIError) *APIError {
	if err == nil {
		return nil
	}
	return apiErr.WithDetails(err.Error())
}

