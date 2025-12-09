package cache

import "errors"

var (
	// ErrCacheMiss é retornado quando uma chave não é encontrada no cache
	ErrCacheMiss = errors.New("cache miss")
	// ErrCacheConnection é retornado quando há erro de conexão com o cache
	ErrCacheConnection = errors.New("cache connection error")
)

