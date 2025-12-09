package api

import (
	"api-go-arquitetura/internal/api/handlers"

	"github.com/gorilla/mux"
)

// NewRouter monta e retorna o router com as rotas registradas pelos handlers
func NewRouter(produtoHandler *handlers.ProdutoHandler, healthCheckHandler *handlers.HealthCheckHandler) *mux.Router {
	router := mux.NewRouter()

	// Rotas versionadas para produtos (v1)
	v1 := router.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/produtos", produtoHandler.GetProdutos).Methods("GET")
	v1.HandleFunc("/produtos/{id}", produtoHandler.GetProduto).Methods("GET")
	v1.HandleFunc("/produtos", produtoHandler.CreateProduto).Methods("POST")
	v1.HandleFunc("/produtos/{id}", produtoHandler.UpdateProduto).Methods("PUT")
	v1.HandleFunc("/produtos/{id}", produtoHandler.PatchProduto).Methods("PATCH")
	v1.HandleFunc("/produtos/{id}", produtoHandler.DeleteProduto).Methods("DELETE")

	// Manter compatibilidade com rotas antigas (redirecionar para v1)
	// Isso permite uma transição suave para o versionamento
	router.HandleFunc("/api/produtos", produtoHandler.GetProdutos).Methods("GET")
	router.HandleFunc("/api/produtos/{id}", produtoHandler.GetProduto).Methods("GET")
	router.HandleFunc("/api/produtos", produtoHandler.CreateProduto).Methods("POST")
	router.HandleFunc("/api/produtos/{id}", produtoHandler.UpdateProduto).Methods("PUT")
	router.HandleFunc("/api/produtos/{id}", produtoHandler.PatchProduto).Methods("PATCH")
	router.HandleFunc("/api/produtos/{id}", produtoHandler.DeleteProduto).Methods("DELETE")

	// Rota de health check (não versionada)
	if healthCheckHandler != nil {
		router.HandleFunc("/health", healthCheckHandler.HealthCheck).Methods("GET")
	}

	return router
}
