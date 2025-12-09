package service

import (
	"context"

	"api-go-arquitetura/internal/dto"
	"api-go-arquitetura/internal/model"
)

// ProdutoService define a interface para operações de produto
type ProdutoService interface {
	Create(ctx context.Context, produto model.Produto) (model.Produto, error)
	FindAll(ctx context.Context) ([]model.Produto, error)
	FindByID(ctx context.Context, id int) (model.Produto, error)
	Update(ctx context.Context, id int, produto model.Produto) (model.Produto, error)
	Patch(ctx context.Context, id int, updates map[string]interface{}) (model.Produto, error)
	Delete(ctx context.Context, id int) error
	// Novos métodos para paginação e filtros
	FindAllPaginated(ctx context.Context, pagination dto.PaginationRequest, filter dto.FilterRequest, sort dto.SortRequest) ([]model.Produto, dto.PaginationResponse, error)
}

