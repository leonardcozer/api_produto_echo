package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "api-go-arquitetura/docs"
	"api-go-arquitetura/internal/api"
	"api-go-arquitetura/internal/api/handlers"
	"api-go-arquitetura/internal/api/middleware"
	"api-go-arquitetura/internal/cache"
	"api-go-arquitetura/internal/config"
	"api-go-arquitetura/internal/database"
	"api-go-arquitetura/internal/logger"
	"api-go-arquitetura/internal/metrics"
	"api-go-arquitetura/internal/repository"
	"api-go-arquitetura/internal/service"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title API Go com Arquitetura
// @version 1.0
// @description Uma API REST completa em Go com suporte aos verbos HTTP
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @basePath /api/v1
// @schemes http
func main() {
	// Carregar configurações
	cfg := config.Load()
	
	// Validar configurações
	if err := cfg.Validate(); err != nil {
		logger.Fatalf("Erro na configuração: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"mongo_uri": cfg.MongoURI,
		"database":  cfg.Database,
		"port":      cfg.Port,
	}).Info("Configurações carregadas")

	// Conectar ao MongoDB com tratamento de erro robusto
	opts := database.ConnectOptions{
		URI:            cfg.MongoURI,
		ConnectTimeout: cfg.ConnectTimeout,
		MaxPoolSize:    cfg.MaxPoolSize,
		MinPoolSize:    cfg.MinPoolSize,
	}
	
	client, err := database.Connect(opts)
	if err != nil {
		logger.WithField("error", err).Fatal("Erro ao conectar ao MongoDB")
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := database.Disconnect(ctx, client); err != nil {
			logger.WithField("error", err).Error("Erro ao fechar conexão com MongoDB")
		}
	}()

	// Obter coleção de produtos com tratamento de erro
	col, err := database.GetCollection(client, cfg.Database, "produtos")
	if err != nil {
		logger.WithField("error", err).Fatal("Erro ao obter coleção")
	}

	// Criar índices otimizados
	ctxIndex, cancelIndex := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelIndex()
	if err := database.CreateIndexes(ctxIndex, client, cfg.Database, "produtos"); err != nil {
		logger.WithField("error", err).Warn("Erro ao criar índices (continuando mesmo assim)")
	}

	// Criar repositório
	prodRepo := repository.NewProdutoRepository(col)

	// Inicializar cache
	var cacheInstance cache.Cache
	if cfg.CacheType == "redis" {
		redisCache, err := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
		if err != nil {
			logger.WithField("error", err).Warn("Erro ao conectar ao Redis, usando cache em memória")
			cacheInstance = cache.NewMemoryCache()
		} else {
			cacheInstance = redisCache
			logger.WithFields(map[string]interface{}{
				"type": "redis",
				"addr": cfg.RedisAddr,
			}).Info("Cache Redis inicializado")
		}
	} else {
		cacheInstance = cache.NewMemoryCache()
		logger.WithField("type", "memory").Info("Cache em memória inicializado")
	}

	// Criar service e injetar o repositório e cache
	prodService := service.NewProdutoServiceWithTTL(prodRepo, cacheInstance, cfg.CacheTTL)

	// Criar handler e injetar o service
	produtoHandler := handlers.NewProdutoHandler(prodService)

	// Criar health check handler com verificação de banco de dados
	healthCheckFunc := func(ctx context.Context) error {
		return database.HealthCheck(ctx, client)
	}
	healthCheckHandler := handlers.NewHealthCheckHandler(healthCheckFunc)

	// Criar router e injetar os handlers
	router := api.NewRouter(produtoHandler, healthCheckHandler)

	// Rota do Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Rota de métricas Prometheus
	router.Handle("/metrics", metrics.GetHandler()).Methods("GET")

	// Configurar CORS
	middleware.SetCORSConfig(&cfg)

	// Aplicar middlewares
	handler := middleware.ApplyMiddlewares(router)

	// Configurar servidor HTTP usando configurações
	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Canal para receber sinais do sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em goroutine
	go func() {
		logger.WithFields(map[string]interface{}{
			"port":    cfg.Port,
			"swagger": "http://localhost" + cfg.Port + "/swagger/index.html",
		}).Info("Servidor iniciando")
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithField("error", err).Fatal("Erro ao iniciar servidor")
		}
	}()

	// Aguardar sinal de interrupção
	<-quit
	logger.Info("Servidor sendo encerrado...")

	// Graceful shutdown usando timeout da config
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithField("error", err).Fatal("Erro ao encerrar servidor")
	}

	logger.Info("Servidor encerrado com sucesso")
	
	// Fazer shutdown do logger (flush final para Loki)
	logger.Shutdown()
}
