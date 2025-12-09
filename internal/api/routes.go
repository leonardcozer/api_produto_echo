package api

import (
	"api-go-arquitetura/internal/api/handlers"

	"github.com/labstack/echo/v4"
)

// NewRouter monta e retorna o router Echo com as rotas registradas pelos handlers
func NewRouter(produtoHandler *handlers.ProdutoHandler, healthCheckHandler *handlers.HealthCheckHandler) *echo.Echo {
	e := echo.New()

	// Rotas versionadas para produtos (v1)
	v1 := e.Group("/api/v1")
	v1.GET("/produtos", produtoHandler.GetProdutos)
	v1.GET("/produtos/:id", produtoHandler.GetProduto)
	v1.POST("/produtos", produtoHandler.CreateProduto)
	v1.PUT("/produtos/:id", produtoHandler.UpdateProduto)
	v1.PATCH("/produtos/:id", produtoHandler.PatchProduto)
	v1.DELETE("/produtos/:id", produtoHandler.DeleteProduto)

	// Manter compatibilidade com rotas antigas (redirecionar para v1)
	// Isso permite uma transição suave para o versionamento
	legacy := e.Group("/api")
	legacy.GET("/produtos", produtoHandler.GetProdutos)
	legacy.GET("/produtos/:id", produtoHandler.GetProduto)
	legacy.POST("/produtos", produtoHandler.CreateProduto)
	legacy.PUT("/produtos/:id", produtoHandler.UpdateProduto)
	legacy.PATCH("/produtos/:id", produtoHandler.PatchProduto)
	legacy.DELETE("/produtos/:id", produtoHandler.DeleteProduto)

	// Rota de health check (não versionada)
	if healthCheckHandler != nil {
		e.GET("/health", healthCheckHandler.HealthCheck)
	}

	return e
}
