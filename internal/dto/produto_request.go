package dto

// CreateProdutoRequest representa os dados necessários para criar um produto
// @Description Dados para criação de um novo produto
type CreateProdutoRequest struct {
	Nome      string  `json:"nome" validate:"required,min=1,max=100" example:"Notebook"`
	Preco     float64 `json:"preco" validate:"required,gt=0" example:"3500.00"`
	Descricao string  `json:"descricao" validate:"max=500" example:"Notebook de alta performance"`
}

// UpdateProdutoRequest representa os dados necessários para atualizar um produto
// @Description Dados para atualização completa de um produto
type UpdateProdutoRequest struct {
	Nome      string  `json:"nome" validate:"required,min=1,max=100" example:"Notebook"`
	Preco     float64 `json:"preco" validate:"required,gt=0" example:"3500.00"`
	Descricao string  `json:"descricao" validate:"max=500" example:"Notebook de alta performance"`
}

// PatchProdutoRequest representa os dados para atualização parcial de um produto
// @Description Dados para atualização parcial de um produto (campos opcionais)
type PatchProdutoRequest struct {
	Nome      *string  `json:"nome,omitempty" validate:"omitempty,min=1,max=100" example:"Notebook"`
	Preco     *float64 `json:"preco,omitempty" validate:"omitempty,gt=0" example:"3500.00"`
	Descricao *string  `json:"descricao,omitempty" validate:"omitempty,max=500" example:"Notebook de alta performance"`
}

