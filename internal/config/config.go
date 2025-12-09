package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config contém todas as configurações da aplicação
type Config struct {
	// MongoDB
	MongoURI      string
	Database      string
	ConnectTimeout time.Duration
	
	// Server
	Port          string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
	ShutdownTimeout time.Duration
	
	// Database Pool
	MaxPoolSize  uint64
	MinPoolSize  uint64
	
	// Observability
	LokiURL string
	LokiJob string
	
	// Cache
	CacheType      string        // "memory" ou "redis"
	CacheTTL       time.Duration // TTL padrão do cache
	RedisAddr      string        // Endereço do Redis (ex: "localhost:6379")
	RedisPassword  string        // Senha do Redis
	RedisDB        int           // Database do Redis
	
	// CORS
	CORSAllowedOrigins []string // Origens permitidas (vazio = todas)
	CORSAllowedMethods []string // Métodos permitidos
	CORSAllowedHeaders []string // Headers permitidos
	CORSCredentials    bool     // Permitir credenciais
}

// Load carrega as configurações da aplicação a partir de variáveis de ambiente
// com valores padrão apropriados
func Load() Config {
	port := getEnv("PORT", "8080")
	// Garantir que a porta tenha o formato correto
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return Config{
		// MongoDB
		MongoURI:       getEnv("MONGO_URI", "mongodb://localhost:27017"),
		Database:        getEnv("MONGO_DB", "api_go"),
		ConnectTimeout:  getDurationEnv("MONGO_CONNECT_TIMEOUT", 10*time.Second),
		
		// Server
		Port:            port,
		ReadTimeout:     getDurationEnv("READ_TIMEOUT", 15*time.Second),
		WriteTimeout:    getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:     getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 30*time.Second),
		
		// Database Pool
		MaxPoolSize: getUint64Env("MONGO_MAX_POOL_SIZE", 100),
		MinPoolSize: getUint64Env("MONGO_MIN_POOL_SIZE", 10),
		
		// Observability
		LokiURL: getEnv("LOKI_URL", ""),
		LokiJob: getEnv("LOKI_JOB", "ARQUITETURA"),
		
		// Cache
		CacheType:     getEnv("CACHE_TYPE", "memory"), // memory ou redis
		CacheTTL:      getDurationEnv("CACHE_TTL", 5*time.Minute),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:        getIntEnv("REDIS_DB", 0),
		
		// CORS
		CORSAllowedOrigins: getStringSliceEnv("CORS_ALLOWED_ORIGINS", []string{"*"}),
		CORSAllowedMethods: getStringSliceEnv("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		CORSAllowedHeaders: getStringSliceEnv("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
		CORSCredentials:    getBoolEnv("CORS_CREDENTIALS", false),
	}
}

// Validate valida as configurações e retorna erro se alguma estiver inválida
func (c *Config) Validate() error {
	if c.MongoURI == "" {
		return fmt.Errorf("MONGO_URI não pode ser vazia")
	}
	if c.Database == "" {
		return fmt.Errorf("MONGO_DB não pode ser vazio")
	}
	if c.Port == "" {
		return fmt.Errorf("PORT não pode ser vazio")
	}
	if c.ConnectTimeout <= 0 {
		return fmt.Errorf("MONGO_CONNECT_TIMEOUT deve ser maior que zero")
	}
	return nil
}

// getEnv obtém uma variável de ambiente ou retorna o valor padrão
func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}

// getDurationEnv obtém uma variável de ambiente como duration ou retorna o valor padrão
func getDurationEnv(key string, def time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return def
}

// getUint64Env obtém uma variável de ambiente como uint64 ou retorna o valor padrão
func getUint64Env(key string, def uint64) uint64 {
	if value := os.Getenv(key); value != "" {
		var result uint64
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return def
}

// getIntEnv obtém uma variável de ambiente como int ou retorna o valor padrão
func getIntEnv(key string, def int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return def
}

// getStringSliceEnv obtém uma variável de ambiente como slice de strings (separado por vírgula) ou retorna o valor padrão
func getStringSliceEnv(key string, def []string) []string {
	if value := os.Getenv(key); value != "" {
		if value == "*" {
			return []string{"*"}
		}
		return strings.Split(value, ",")
	}
	return def
}

// getBoolEnv obtém uma variável de ambiente como bool ou retorna o valor padrão
func getBoolEnv(key string, def bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true" || value == "1"
	}
	return def
}
