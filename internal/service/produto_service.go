package service

import (
	"context"
	"time"

	"api-go-arquitetura/internal/cache"
	"api-go-arquitetura/internal/dto"
	"api-go-arquitetura/internal/errors"
	"api-go-arquitetura/internal/logger"
	"api-go-arquitetura/internal/metrics"
	"api-go-arquitetura/internal/model"
	"api-go-arquitetura/internal/repository"
)

// produtoService implementa a lógica de negócio para produtos
type produtoService struct {
	repo  repository.ProdutoRepository
	cache cache.Cache
	ttl   time.Duration
}

// NewProdutoService cria uma nova instância do ProdutoService
func NewProdutoService(repo repository.ProdutoRepository, cache cache.Cache) ProdutoService {
	// TTL padrão de 5 minutos para cache
	ttl := 5 * time.Minute
	return &produtoService{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

// NewProdutoServiceWithTTL cria uma nova instância do ProdutoService com TTL customizado
func NewProdutoServiceWithTTL(repo repository.ProdutoRepository, cache cache.Cache, ttl time.Duration) ProdutoService {
	return &produtoService{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

// Create cria um novo produto
func (s *produtoService) Create(ctx context.Context, produto model.Produto) (model.Produto, error) {
	// Validações de negócio
	if produto.Nome == "" {
		return model.Produto{}, errors.ErrNomeObrigatorio
	}
	if produto.Preco <= 0 {
		return model.Produto{}, errors.ErrPrecoInvalido
	}

	result, err := s.repo.Create(ctx, produto)
	if err != nil {
		return model.Produto{}, errors.WrapError(err, errors.ErrDatabase)
	}

	// Invalidar cache de listas (novo produto adicionado)
	if s.cache != nil {
		// Em produção, seria melhor usar padrões de chave ou tags do Redis
		// Por enquanto, o cache será invalidado naturalmente pelo TTL
		logger.Debug("Cache de listas será invalidado pelo TTL após criação de produto")
	}

	return result, nil
}

// FindAll retorna todos os produtos
func (s *produtoService) FindAll(ctx context.Context) ([]model.Produto, error) {
	result, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return result, nil
}

// FindByID retorna um produto pelo ID
func (s *produtoService) FindByID(ctx context.Context, id int) (model.Produto, error) {
	if id <= 0 {
		return model.Produto{}, errors.ErrInvalidID
	}

	// Tentar buscar do cache primeiro
	if s.cache != nil {
		cacheKey := cache.GenerateProdutoKey(id)
		start := time.Now()
		cachedData, err := s.cache.Get(ctx, cacheKey)
		duration := time.Since(start)
		
		if err == nil {
			produto, err := cache.DecodeProduto(cachedData)
			if err == nil {
				metrics.RecordCacheHit("get", duration)
				logger.WithFields(map[string]interface{}{
					"id":        id,
					"cache_key": cacheKey,
				}).Debug("Cache hit para produto")
				return produto, nil
			}
			metrics.RecordCacheError("get", duration)
		} else {
			metrics.RecordCacheMiss("get", duration)
		}
	}

	// Cache miss ou erro - buscar do banco
	result, err := s.repo.FindByID(ctx, id)
	if err != nil {
		// Verificar se é erro de "not found" do repository
		if err.Error() == "not found" {
			return model.Produto{}, errors.ErrProdutoNotFound
		}
		return model.Produto{}, errors.WrapError(err, errors.ErrDatabase)
	}

	// Armazenar no cache
	if s.cache != nil {
		cacheKey := cache.GenerateProdutoKey(id)
		cachedData, err := cache.EncodeProduto(result)
		if err == nil {
			start := time.Now()
			if err := s.cache.Set(ctx, cacheKey, cachedData, s.ttl); err != nil {
				metrics.RecordCacheError("set", time.Since(start))
				logger.WithField("error", err).Warn("Erro ao armazenar produto no cache")
			} else {
				metrics.RecordCacheOperation("set", "success", time.Since(start))
			}
		}
	}

	return result, nil
}

// Update atualiza um produto completamente
func (s *produtoService) Update(ctx context.Context, id int, produto model.Produto) (model.Produto, error) {
	if id <= 0 {
		return model.Produto{}, errors.ErrInvalidID
	}
	if produto.Nome == "" {
		return model.Produto{}, errors.ErrNomeObrigatorio
	}
	if produto.Preco <= 0 {
		return model.Produto{}, errors.ErrPrecoInvalido
	}

	result, err := s.repo.Update(ctx, id, produto)
	if err != nil {
		if err.Error() == "not found" {
			return model.Produto{}, errors.ErrProdutoNotFound
		}
		return model.Produto{}, errors.WrapError(err, errors.ErrDatabase)
	}

	// Invalidar cache do produto atualizado
	if s.cache != nil {
		cacheKey := cache.GenerateProdutoKey(id)
		start := time.Now()
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			metrics.RecordCacheError("delete", time.Since(start))
			logger.WithField("error", err).Warn("Erro ao invalidar cache do produto")
		} else {
			metrics.RecordCacheOperation("delete", "success", time.Since(start))
		}
		// Invalidar cache de listas também
		logger.Debug("Cache invalidado após atualização de produto")
	}

	return result, nil
}

// Patch atualiza um produto parcialmente
func (s *produtoService) Patch(ctx context.Context, id int, updates map[string]interface{}) (model.Produto, error) {
	if id <= 0 {
		return model.Produto{}, errors.ErrInvalidID
	}

	// Validações específicas para campos que podem ser atualizados
	if nome, ok := updates["nome"].(string); ok && nome == "" {
		return model.Produto{}, errors.ErrNomeObrigatorio
	}
	if preco, ok := updates["preco"].(float64); ok && preco <= 0 {
		return model.Produto{}, errors.ErrPrecoInvalido
	}

	result, err := s.repo.Patch(ctx, id, updates)
	if err != nil {
		if err.Error() == "not found" {
			return model.Produto{}, errors.ErrProdutoNotFound
		}
		return model.Produto{}, errors.WrapError(err, errors.ErrDatabase)
	}

	// Invalidar cache do produto atualizado
	if s.cache != nil {
		cacheKey := cache.GenerateProdutoKey(id)
		start := time.Now()
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			metrics.RecordCacheError("delete", time.Since(start))
			logger.WithField("error", err).Warn("Erro ao invalidar cache do produto")
		} else {
			metrics.RecordCacheOperation("delete", "success", time.Since(start))
		}
		// Invalidar cache de listas também
		logger.Debug("Cache invalidado após patch de produto")
	}

	return result, nil
}

// Delete remove um produto
func (s *produtoService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.ErrInvalidID
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err.Error() == "not found" {
			return errors.ErrProdutoNotFound
		}
		return errors.WrapError(err, errors.ErrDatabase)
	}

	// Invalidar cache do produto deletado
	if s.cache != nil {
		cacheKey := cache.GenerateProdutoKey(id)
		start := time.Now()
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			metrics.RecordCacheError("delete", time.Since(start))
			logger.WithField("error", err).Warn("Erro ao invalidar cache do produto")
		} else {
			metrics.RecordCacheOperation("delete", "success", time.Since(start))
		}
		// Invalidar cache de listas também
		logger.Debug("Cache invalidado após deleção de produto")
	}

	return nil
}

// FindAllPaginated retorna produtos paginados com filtros e ordenação
func (s *produtoService) FindAllPaginated(ctx context.Context, pagination dto.PaginationRequest, filter dto.FilterRequest, sort dto.SortRequest) ([]model.Produto, dto.PaginationResponse, error) {
	// Validar paginação
	pagination.Validate()

	// Validar ordenação
	if err := sort.Validate(); err != nil {
		return nil, dto.PaginationResponse{}, errors.ErrInvalidInput.WithDetails(err.Error())
	}

	// Converter filtro para MongoDB
	mongoFilter := filter.ToMongoFilter()

	// Converter ordenação para MongoDB
	mongoSort := sort.ToMongoSort()

	// Gerar chave de cache para a lista
	cacheKey := cache.GenerateProdutosListKey(pagination.Page, pagination.PageSize, mongoFilter)

	// Tentar buscar do cache primeiro
	if s.cache != nil {
		start := time.Now()
		cachedData, err := s.cache.Get(ctx, cacheKey)
		duration := time.Since(start)
		
		if err == nil {
			// Cache hit
			var cachedResult struct {
				Produtos []model.Produto
				Total    int64
			}
			if err := cache.Decode(cachedData, &cachedResult); err == nil {
				metrics.RecordCacheHit("get_list", duration)
				logger.WithFields(map[string]interface{}{
					"cache_key": cacheKey,
					"page":       pagination.Page,
				}).Debug("Cache hit para lista de produtos")
				
				paginationResp := dto.NewPaginationResponse(pagination.Page, pagination.PageSize, int(cachedResult.Total))
				return cachedResult.Produtos, paginationResp, nil
			}
			metrics.RecordCacheError("get_list", duration)
		} else {
			metrics.RecordCacheMiss("get_list", duration)
		}
	}

	// Contar total de documentos
	totalItems, err := s.repo.Count(ctx, mongoFilter)
	if err != nil {
		return nil, dto.PaginationResponse{}, errors.WrapError(err, errors.ErrDatabase)
	}

	// Buscar produtos paginados
	produtos, err := s.repo.FindAllPaginated(ctx, pagination.GetSkip(), pagination.GetLimit(), mongoFilter, mongoSort)
	if err != nil {
		return nil, dto.PaginationResponse{}, errors.WrapError(err, errors.ErrDatabase)
	}

	// Armazenar no cache
	if s.cache != nil {
		cachedResult := struct {
			Produtos []model.Produto
			Total    int64
		}{
			Produtos: produtos,
			Total:    totalItems,
		}
		if cachedData, err := cache.Encode(cachedResult); err == nil {
			start := time.Now()
			if err := s.cache.Set(ctx, cacheKey, cachedData, s.ttl); err != nil {
				metrics.RecordCacheError("set_list", time.Since(start))
				logger.WithField("error", err).Warn("Erro ao armazenar lista no cache")
			} else {
				metrics.RecordCacheOperation("set_list", "success", time.Since(start))
			}
		}
	}

	// Criar resposta de paginação
	paginationResp := dto.NewPaginationResponse(pagination.Page, pagination.PageSize, int(totalItems))

	return produtos, paginationResp, nil
}
