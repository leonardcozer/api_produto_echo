package cache

import (
	"context"
	"sync"
	"time"
)

// memoryCache implementa Cache usando memória local
type memoryCache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
}

type cacheItem struct {
	value      []byte
	expiration time.Time
}

// NewMemoryCache cria uma nova instância de cache em memória
func NewMemoryCache() Cache {
	c := &memoryCache{
		items: make(map[string]*cacheItem),
	}
	// Iniciar goroutine para limpar itens expirados
	go c.cleanup()
	return c
}

// Get recupera um valor do cache
func (c *memoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, ErrCacheMiss
	}

	// Verificar se o item expirou
	if time.Now().After(item.expiration) {
		// Remover item expirado (em modo de escrita)
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		c.mu.RLock()
		return nil, ErrCacheMiss
	}

	// Retornar cópia do valor para evitar race conditions
	result := make([]byte, len(item.value))
	copy(result, item.value)
	return result, nil
}

// Set armazena um valor no cache com TTL
func (c *memoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(ttl)
	c.items[key] = &cacheItem{
		value:      value,
		expiration: expiration,
	}
	return nil
}

// Delete remove um valor do cache
func (c *memoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

// Clear limpa todo o cache
func (c *memoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
	return nil
}

// Exists verifica se uma chave existe no cache
func (c *memoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return false, nil
	}

	// Verificar se expirou
	if time.Now().After(item.expiration) {
		return false, nil
	}

	return true, nil
}

// cleanup remove itens expirados periodicamente
func (c *memoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.expiration) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

