package repository

import (
	"context"

	"api-go-arquitetura/internal/model"

	"go.mongodb.org/mongo-driver/bson"
)

// ProdutoRepository define a interface para operações de produto no repositório
type ProdutoRepository interface {
	Create(ctx context.Context, produto model.Produto) (model.Produto, error)
	FindAll(ctx context.Context) ([]model.Produto, error)
	FindByID(ctx context.Context, id int) (model.Produto, error)
	Update(ctx context.Context, id int, produto model.Produto) (model.Produto, error)
	Patch(ctx context.Context, id int, updates map[string]interface{}) (model.Produto, error)
	Delete(ctx context.Context, id int) error
	// Novos métodos para paginação e filtros
	FindAllPaginated(ctx context.Context, skip, limit int64, filter map[string]interface{}, sort bson.D) ([]model.Produto, error)
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}

