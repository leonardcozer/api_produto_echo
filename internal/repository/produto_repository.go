package repository

import (
	"context"
	"errors"
	"time"

	"api-go-arquitetura/internal/database"
	"api-go-arquitetura/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoProdutoRepository implementa ProdutoRepository usando MongoDB
type mongoProdutoRepository struct {
	Collection *mongo.Collection
}

// NewProdutoRepository cria uma nova instância do ProdutoRepository
func NewProdutoRepository(col *mongo.Collection) ProdutoRepository {
	return &mongoProdutoRepository{Collection: col}
}

func (r *mongoProdutoRepository) getNextID(ctx context.Context) (int, error) {
	opts := options.FindOne().SetSort(bson.D{{"id", -1}})
	var p model.Produto
	err := r.Collection.FindOne(ctx, bson.M{}, opts).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil
		}
		return 0, err
	}
	return p.ID + 1, nil
}

func (r *mongoProdutoRepository) Create(ctx context.Context, produto model.Produto) (model.Produto, error) {
	// Usar retry logic para operação crítica
	retryOpts := database.DefaultRetryOptions()
	result, err := database.RetryWithResult(ctx, func() (model.Produto, error) {
		id, err := r.getNextID(ctx)
		if err != nil {
			return model.Produto{}, err
		}
		produto.ID = id
		produto.BeforeCreate() // Inicializar timestamps
		_, err = r.Collection.InsertOne(ctx, produto)
		if err != nil {
			return model.Produto{}, err
		}
		return produto, nil
	}, retryOpts)
	
	return result, err
}

func (r *mongoProdutoRepository) FindAll(ctx context.Context) ([]model.Produto, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var produtos []model.Produto
	if err = cursor.All(ctx, &produtos); err != nil {
		return nil, err
	}
	return produtos, nil
}

func (r *mongoProdutoRepository) FindByID(ctx context.Context, id int) (model.Produto, error) {
	// Filtrar produtos deletados (soft delete)
	filter := bson.M{
		"id":        id,
		"deleted_at": bson.M{"$exists": false},
	}
	var produto model.Produto
	err := r.Collection.FindOne(ctx, filter).Decode(&produto)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Produto{}, errors.New("not found")
		}
		return model.Produto{}, err
	}
	return produto, nil
}

func (r *mongoProdutoRepository) Update(ctx context.Context, id int, produto model.Produto) (model.Produto, error) {
	// Usar retry logic para operação crítica
	retryOpts := database.DefaultRetryOptions()
	result, err := database.RetryWithResult(ctx, func() (model.Produto, error) {
		produto.ID = id
		produto.BeforeUpdate() // Atualizar timestamp
		
		// Buscar produto existente para preservar CreatedAt e verificar se não está deletado
		filter := bson.M{
			"id":        id,
			"deleted_at": bson.M{"$exists": false},
		}
		var existing model.Produto
		err := r.Collection.FindOne(ctx, filter).Decode(&existing)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return model.Produto{}, errors.New("not found")
			}
			return model.Produto{}, err
		}
		produto.CreatedAt = existing.CreatedAt // Preservar CreatedAt
		produto.DeletedAt = existing.DeletedAt // Preservar DeletedAt (soft delete)
		
		res, err := r.Collection.ReplaceOne(ctx, filter, produto)
		if err != nil {
			return model.Produto{}, err
		}
		if res.MatchedCount == 0 {
			return model.Produto{}, errors.New("not found")
		}
		return produto, nil
	}, retryOpts)
	
	return result, err
}

func (r *mongoProdutoRepository) Patch(ctx context.Context, id int, updates map[string]interface{}) (model.Produto, error) {
	// Adicionar updated_at automaticamente
	updates["updated_at"] = time.Now()
	
	// Filtrar produtos deletados (soft delete)
	filter := bson.M{
		"id":        id,
		"deleted_at": bson.M{"$exists": false},
	}
	
	update := bson.M{"$set": updates}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated model.Produto
	err := r.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Produto{}, errors.New("not found")
		}
		return model.Produto{}, err
	}
	return updated, nil
}

func (r *mongoProdutoRepository) Delete(ctx context.Context, id int) error {
	// Soft delete: marcar como deletado ao invés de remover
	retryOpts := database.DefaultRetryOptions()
	err := database.Retry(ctx, func() error {
		now := time.Now()
		update := bson.M{
			"$set": bson.M{
				"deleted_at": now,
				"updated_at": now,
			},
		}
		res, err := r.Collection.UpdateOne(ctx, bson.M{"id": id, "deleted_at": bson.M{"$exists": false}}, update)
		if err != nil {
			return err
		}
		if res.MatchedCount == 0 {
			return errors.New("not found")
		}
		return nil
	}, retryOpts)
	
	return err
}

// FindAllPaginated retorna produtos paginados com filtros e ordenação
func (r *mongoProdutoRepository) FindAllPaginated(ctx context.Context, skip, limit int64, filter map[string]interface{}, sort bson.D) ([]model.Produto, error) {
	// Converter filter para bson.M
	mongoFilter := bson.M{}
	if filter != nil {
		mongoFilter = bson.M(filter)
	}
	// Filtrar produtos deletados (soft delete)
	mongoFilter["deleted_at"] = bson.M{"$exists": false}

	// Se sort estiver vazio, usar ordenação padrão por ID
	if len(sort) == 0 {
		sort = bson.D{{Key: "id", Value: 1}}
	}

	// Opções de paginação
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(sort)

	cursor, err := r.Collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var produtos []model.Produto
	if err = cursor.All(ctx, &produtos); err != nil {
		return nil, err
	}
	return produtos, nil
}

// Count retorna o total de documentos que correspondem ao filtro
func (r *mongoProdutoRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	// Converter filter para bson.M
	mongoFilter := bson.M{}
	if filter != nil {
		mongoFilter = bson.M(filter)
	}
	// Filtrar produtos deletados (soft delete)
	mongoFilter["deleted_at"] = bson.M{"$exists": false}

	count, err := r.Collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
