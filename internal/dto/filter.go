package dto

// FilterRequest representa os filtros de busca
type FilterRequest struct {
	Nome      *string  `json:"nome,omitempty"`      // Busca por nome (contém)
	PrecoMin  *float64 `json:"precoMin,omitempty"`  // Preço mínimo
	PrecoMax  *float64 `json:"precoMax,omitempty"`  // Preço máximo
	Descricao *string  `json:"descricao,omitempty"` // Busca por descrição (contém)
}

// ToMongoFilter converte FilterRequest para filtro MongoDB
func (f *FilterRequest) ToMongoFilter() map[string]interface{} {
	filter := make(map[string]interface{})

	if f.Nome != nil && *f.Nome != "" {
		filter["nome"] = map[string]interface{}{
			"$regex":   *f.Nome,
			"$options": "i", // Case insensitive
		}
	}

	if f.Descricao != nil && *f.Descricao != "" {
		filter["descricao"] = map[string]interface{}{
			"$regex":   *f.Descricao,
			"$options": "i", // Case insensitive
		}
	}

	// Filtro de preço
	precoFilter := make(map[string]interface{})
	if f.PrecoMin != nil {
		precoFilter["$gte"] = *f.PrecoMin
	}
	if f.PrecoMax != nil {
		precoFilter["$lte"] = *f.PrecoMax
	}
	if len(precoFilter) > 0 {
		filter["preco"] = precoFilter
	}

	return filter
}

// IsEmpty verifica se o filtro está vazio
func (f *FilterRequest) IsEmpty() bool {
	return (f.Nome == nil || *f.Nome == "") &&
		(f.PrecoMin == nil) &&
		(f.PrecoMax == nil) &&
		(f.Descricao == nil || *f.Descricao == "")
}

