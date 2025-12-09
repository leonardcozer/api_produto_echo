package cache

import (
	"context"
	"fmt"
	"time"
)

// Cache define a interface para operações de cache
type Cache interface {
	// Get recupera um valor do cache
	Get(ctx context.Context, key string) ([]byte, error)
	// Set armazena um valor no cache com TTL
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	// Delete remove um valor do cache
	Delete(ctx context.Context, key string) error
	// Clear limpa todo o cache
	Clear(ctx context.Context) error
	// Exists verifica se uma chave existe no cache
	Exists(ctx context.Context, key string) (bool, error)
}

// KeyGenerator gera chaves de cache de forma consistente
type KeyGenerator struct {
	prefix string
}

// NewKeyGenerator cria um novo gerador de chaves
func NewKeyGenerator(prefix string) *KeyGenerator {
	return &KeyGenerator{prefix: prefix}
}

// Generate gera uma chave de cache
func (kg *KeyGenerator) Generate(parts ...string) string {
	key := kg.prefix
	for _, part := range parts {
		if part != "" {
			key += ":" + part
		}
	}
	return key
}

// ProdutoKeyGenerator gera chaves específicas para produtos
var ProdutoKeyGenerator = NewKeyGenerator("produto")

// GenerateProdutoKey gera uma chave de cache para um produto
func GenerateProdutoKey(id int) string {
	return ProdutoKeyGenerator.Generate("id", fmt.Sprintf("%d", id))
}

// GenerateProdutosListKey gera uma chave de cache para lista de produtos
func GenerateProdutosListKey(page, pageSize int, filters map[string]interface{}) string {
	key := ProdutoKeyGenerator.Generate("list")
	if page > 0 {
		key += ":page:" + fmt.Sprintf("%d", page)
	}
	if pageSize > 0 {
		key += ":size:" + fmt.Sprintf("%d", pageSize)
	}
	// Adicionar filtros à chave se existirem
	if filters != nil && len(filters) > 0 {
		// Simplificado: em produção, seria melhor usar hash dos filtros
		for k, v := range filters {
			key += ":" + k + ":" + fmt.Sprintf("%v", v)
		}
	}
	return key
}

// InvalidateListCache invalida todas as listas em cache
// Em produção, seria melhor usar padrões de chave ou tags
func InvalidateListCache(ctx context.Context, cache Cache) error {
	// Por enquanto, não implementamos limpeza completa
	// O cache será invalidado naturalmente pelo TTL
	// Em produção, usaríamos Redis com padrões de chave ou tags
	return nil
}

