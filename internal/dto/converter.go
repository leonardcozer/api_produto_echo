package dto

import "api-go-arquitetura/internal/model"

// ToModel converte CreateProdutoRequest para model.Produto
func (r *CreateProdutoRequest) ToModel() model.Produto {
	return model.Produto{
		Nome:      r.Nome,
		Preco:     r.Preco,
		Descricao: r.Descricao,
	}
}

// ToModel converte UpdateProdutoRequest para model.Produto
func (r *UpdateProdutoRequest) ToModel() model.Produto {
	return model.Produto{
		Nome:      r.Nome,
		Preco:     r.Preco,
		Descricao: r.Descricao,
	}
}

// ToMap converte PatchProdutoRequest para map[string]interface{}
func (r *PatchProdutoRequest) ToMap() map[string]interface{} {
	updates := make(map[string]interface{})
	
	if r.Nome != nil {
		updates["nome"] = *r.Nome
	}
	if r.Preco != nil {
		updates["preco"] = *r.Preco
	}
	if r.Descricao != nil {
		updates["descricao"] = *r.Descricao
	}
	
	return updates
}

// FromModel converte model.Produto para ProdutoResponse
func FromModel(p model.Produto) ProdutoResponse {
	return ProdutoResponse{
		ID:        p.ID,
		Nome:      p.Nome,
		Preco:     p.Preco,
		Descricao: p.Descricao,
	}
}

// FromModelList converte []model.Produto para []ProdutoResponse
func FromModelList(produtos []model.Produto) []ProdutoResponse {
	responses := make([]ProdutoResponse, len(produtos))
	for i, p := range produtos {
		responses[i] = FromModel(p)
	}
	return responses
}

// ToProdutoListResponse converte []model.Produto para ProdutoListResponse
func ToProdutoListResponse(produtos []model.Produto) ProdutoListResponse {
	return ProdutoListResponse{
		Produtos: FromModelList(produtos),
		Total:    len(produtos),
	}
}

