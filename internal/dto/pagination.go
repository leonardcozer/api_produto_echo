package dto

// PaginationRequest representa os parâmetros de paginação
type PaginationRequest struct {
	Page     int `json:"page"`     // Página atual (começa em 1)
	PageSize int `json:"pageSize"` // Tamanho da página
}

// PaginationResponse representa os metadados de paginação na resposta
type PaginationResponse struct {
	Page       int `json:"page"`       // Página atual
	PageSize   int `json:"pageSize"`   // Tamanho da página
	TotalPages int `json:"totalPages"` // Total de páginas
	TotalItems int `json:"totalItems"`  // Total de itens
	HasNext    bool `json:"hasNext"`   // Tem próxima página
	HasPrev    bool `json:"hasPrev"`   // Tem página anterior
}

// Validate valida os parâmetros de paginação e aplica valores padrão
func (p *PaginationRequest) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10 // Valor padrão
	}
	if p.PageSize > 100 {
		p.PageSize = 100 // Limite máximo
	}
}

// GetSkip calcula o número de documentos a pular
func (p *PaginationRequest) GetSkip() int64 {
	return int64((p.Page - 1) * p.PageSize)
}

// GetLimit retorna o limite de documentos
func (p *PaginationRequest) GetLimit() int64 {
	return int64(p.PageSize)
}

// NewPaginationResponse cria uma resposta de paginação
func NewPaginationResponse(page, pageSize, totalItems int) PaginationResponse {
	totalPages := (totalItems + pageSize - 1) / pageSize // Arredondamento para cima
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

