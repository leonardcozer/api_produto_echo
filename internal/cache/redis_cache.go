package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCache implementa Cache usando Redis
type redisCache struct {
	client *redis.Client
}

// NewRedisCache cria uma nova instância de cache Redis
func NewRedisCache(addr string, password string, db int) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Verificar conexão
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &redisCache{client: client}, nil
}

// Get recupera um valor do cache
func (c *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Set armazena um valor no cache com TTL
func (c *redisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// Delete remove um valor do cache
func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Clear limpa todo o cache (apenas chaves com prefixo, se configurado)
func (c *redisCache) Clear(ctx context.Context) error {
	// Em produção, seria melhor limpar apenas chaves com prefixo específico
	// Por segurança, não implementamos limpeza completa aqui
	// Use Delete para remover chaves específicas
	return nil
}

// Exists verifica se uma chave existe no cache
func (c *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

