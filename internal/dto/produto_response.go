package dto

// ProdutoResponse representa a resposta de um produto
// @Description Resposta com dados do produto
type ProdutoResponse struct {
	ID        int     `json:"id" example:"1"`
	Nome      string  `json:"nome" example:"Notebook"`
	Preco     float64 `json:"preco" example:"3500.00"`
	Descricao string  `json:"descricao" example:"Notebook de alta performance"`
}

// ProdutoListResponse representa uma lista de produtos
// @Description Lista de produtos
type ProdutoListResponse struct {
	Produtos []ProdutoResponse `json:"produtos"`
	Total    int               `json:"total"`
}

