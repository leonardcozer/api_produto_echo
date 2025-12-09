package dto

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// SortRequest representa os parâmetros de ordenação
type SortRequest struct {
	Field string `json:"field" query:"sort" example:"preco"` // Campo para ordenar (ex: "preco", "nome", "created_at")
	Order string `json:"order" query:"order" example:"asc"`   // Ordem: "asc" ou "desc" (padrão: "asc")
}

// Validate valida os parâmetros de ordenação
func (s *SortRequest) Validate() error {
	if s.Field == "" {
		return nil // Sem ordenação
	}

	// Campos permitidos para ordenação
	allowedFields := map[string]bool{
		"id":         true,
		"nome":       true,
		"preco":      true,
		"descricao":  true,
		"created_at": true,
		"updated_at": true,
	}

	if !allowedFields[s.Field] {
		return fmt.Errorf("campo de ordenação inválido: %s. Campos permitidos: id, nome, preco, descricao, created_at, updated_at", s.Field)
	}

	// Normalizar ordem
	s.Order = strings.ToLower(s.Order)
	if s.Order != "asc" && s.Order != "desc" {
		s.Order = "asc" // Padrão
	}

	return nil
}

// ToMongoSort converte SortRequest para bson.D (formato de ordenação do MongoDB)
func (s *SortRequest) ToMongoSort() bson.D {
	if s.Field == "" {
		// Ordenação padrão por ID
		return bson.D{{Key: "id", Value: 1}}
	}

	order := 1 // asc
	if s.Order == "desc" {
		order = -1
	}

	return bson.D{{Key: s.Field, Value: order}}
}

// GetSortFromQuery extrai parâmetros de ordenação da query string
// Formato esperado: ?sort=preco&order=desc ou ?sort=preco:desc
func GetSortFromQuery(sortParam, orderParam string) SortRequest {
	sort := SortRequest{}

	// Se sort contém ":", tratar como formato "campo:ordem"
	if strings.Contains(sortParam, ":") {
		parts := strings.Split(sortParam, ":")
		if len(parts) == 2 {
			sort.Field = parts[0]
			sort.Order = parts[1]
		}
	} else {
		sort.Field = sortParam
		sort.Order = orderParam
	}

	// Se order não foi especificado, usar padrão
	if sort.Order == "" {
		sort.Order = "asc"
	}

	return sort
}

