package dto

// PaginatedProdutoListResponse representa uma resposta paginada de produtos
type PaginatedProdutoListResponse struct {
	Produtos   []ProdutoResponse  `json:"produtos"`
	Pagination PaginationResponse `json:"pagination"`
}

// ToPaginatedResponse converte lista de produtos com paginação
func ToPaginatedResponse(produtos []ProdutoResponse, pagination PaginationResponse) PaginatedProdutoListResponse {
	return PaginatedProdutoListResponse{
		Produtos:   produtos,
		Pagination: pagination,
	}
}

