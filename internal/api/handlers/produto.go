package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"api-go-arquitetura/internal/dto"
	"api-go-arquitetura/internal/errors"
	"api-go-arquitetura/internal/service"
	"api-go-arquitetura/internal/utils"
	"api-go-arquitetura/internal/validator"
)

// ProdutoHandler gerencia os handlers de produto
type ProdutoHandler struct {
	service service.ProdutoService
}

// NewProdutoHandler cria uma nova instância do ProdutoHandler
func NewProdutoHandler(svc service.ProdutoService) *ProdutoHandler {
	return &ProdutoHandler{
		service: svc,
	}
}

// GetProdutos lista todos os produtos (com suporte a paginação, filtros e ordenação)
// @Summary Lista produtos com paginação, filtros e ordenação
// @Description Retorna uma lista paginada de produtos com suporte a filtros e ordenação
// @Tags produtos
// @Accept json
// @Produce json
// @Param page query int false "Número da página (padrão: 1)" default(1)
// @Param pageSize query int false "Tamanho da página (padrão: 10, máximo: 100)" default(10)
// @Param nome query string false "Filtro por nome (busca parcial, case-insensitive)"
// @Param precoMin query number false "Preço mínimo"
// @Param precoMax query number false "Preço máximo"
// @Param descricao query string false "Filtro por descrição (busca parcial, case-insensitive)"
// @Param sort query string false "Campo para ordenação (id, nome, preco, descricao, created_at, updated_at)" default(id)
// @Param order query string false "Ordem de ordenação (asc, desc)" default(asc)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} errors.APIError
// @Failure 500 {object} errors.APIError
// @Router /api/v1/produtos [get]
// GET /api/v1/produtos?page=1&pageSize=10&nome=notebook&precoMin=1000&precoMax=5000&sort=preco&order=desc
func (h *ProdutoHandler) GetProdutos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Parse de parâmetros de paginação
	pagination := dto.PaginationRequest{
		Page:     getIntQuery(r, "page", 1),
		PageSize: getIntQuery(r, "pageSize", 10),
	}

	// Parse de filtros
	filter := dto.FilterRequest{
		Nome:      getStringQuery(r, "nome"),
		PrecoMin:  getFloatQuery(r, "precoMin"),
		PrecoMax:  getFloatQuery(r, "precoMax"),
		Descricao: getStringQuery(r, "descricao"),
	}

	// Parse de ordenação
	sort := dto.GetSortFromQuery(
		r.URL.Query().Get("sort"),
		r.URL.Query().Get("order"),
	)
	
	// Validar ordenação
	if validationErrors := validator.Validate(&sort); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(w, validationErrors)
		return
	}

	// Se não há filtros e paginação padrão, usar método antigo para compatibilidade
	if filter.IsEmpty() && pagination.Page == 1 && pagination.PageSize == 10 && sort.Field == "" {
		// Verificar se há parâmetros de query explícitos
		if r.URL.Query().Get("page") == "" && r.URL.Query().Get("pageSize") == "" && r.URL.Query().Get("sort") == "" {
			// Usar método antigo (sem paginação)
			produtos, err := h.service.FindAll(ctx)
			if err != nil {
				utils.ErrorResponse(w, errors.WrapError(err, errors.ErrDatabase))
				return
			}
			response := dto.ToProdutoListResponse(produtos)
			utils.SuccessResponse(w, http.StatusOK, response)
			return
		}
	}

	// Usar método paginado
	produtos, paginationResp, err := h.service.FindAllPaginated(ctx, pagination, filter, sort)
	if err != nil {
		utils.ErrorResponse(w, errors.WrapError(err, errors.ErrDatabase))
		return
	}

	// Converter models para DTOs
	produtosDTO := dto.FromModelList(produtos)
	response := dto.ToPaginatedResponse(produtosDTO, paginationResp)

	utils.SuccessResponse(w, http.StatusOK, response)
}

// getIntQuery obtém um parâmetro de query como int
func getIntQuery(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	var result int
	if _, err := fmt.Sscanf(value, "%d", &result); err != nil {
		return defaultValue
	}
	return result
}

// getStringQuery obtém um parâmetro de query como string (retorna nil se vazio)
func getStringQuery(r *http.Request, key string) *string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	return &value
}

// getFloatQuery obtém um parâmetro de query como float64 (retorna nil se vazio)
func getFloatQuery(r *http.Request, key string) *float64 {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	var result float64
	if _, err := fmt.Sscanf(value, "%f", &result); err != nil {
		return nil
	}
	return &result
}

// GetProduto obtém um produto por ID
// GET /api/produtos/{id}
func (h *ProdutoHandler) GetProduto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, errors.ErrInvalidID)
		return
	}

	ctx := r.Context()
	produto, err := h.service.FindByID(ctx, id)
	if err != nil {
		if errors.IsAPIError(err) {
			utils.ErrorResponse(w, err)
		} else {
			utils.ErrorResponse(w, errors.ErrProdutoNotFound)
		}
		return
	}

	// Converter model para DTO
	response := dto.FromModel(produto)

	utils.SuccessResponse(w, http.StatusOK, response)
}

// CreateProduto cria um novo produto
// POST /api/produtos
func (h *ProdutoHandler) CreateProduto(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateProdutoRequest
	
	// Decodificar JSON
	if err := utils.DecodeJSON(r.Body, &request); err != nil {
		utils.BadRequestResponse(w, "Erro ao decodificar JSON: "+err.Error())
		return
	}

	// Validar DTO
	if validationErrors := validator.Validate(&request); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(w, validationErrors)
		return
	}

	// Converter DTO para model
	produto := request.ToModel()

	ctx := r.Context()
	created, err := h.service.Create(ctx, produto)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	// Converter model para DTO de resposta
	response := dto.FromModel(created)

	utils.SuccessResponse(w, http.StatusCreated, response)
}

// UpdateProduto atualiza um produto completamente
// PUT /api/produtos/{id}
func (h *ProdutoHandler) UpdateProduto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, errors.ErrInvalidID)
		return
	}

	var request dto.UpdateProdutoRequest
	
	// Decodificar JSON
	if err := utils.DecodeJSON(r.Body, &request); err != nil {
		utils.BadRequestResponse(w, "Erro ao decodificar JSON: "+err.Error())
		return
	}

	// Validar DTO
	if validationErrors := validator.Validate(&request); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(w, validationErrors)
		return
	}

	// Converter DTO para model
	produto := request.ToModel()

	ctx := r.Context()
	updated, err := h.service.Update(ctx, id, produto)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	// Converter model para DTO de resposta
	response := dto.FromModel(updated)

	utils.SuccessResponse(w, http.StatusOK, response)
}

// PatchProduto atualiza um produto parcialmente
// PATCH /api/produtos/{id}
func (h *ProdutoHandler) PatchProduto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, errors.ErrInvalidID)
		return
	}

	var request dto.PatchProdutoRequest
	
	// Decodificar JSON
	if err := utils.DecodeJSON(r.Body, &request); err != nil {
		utils.BadRequestResponse(w, "Erro ao decodificar JSON: "+err.Error())
		return
	}

	// Validar DTO (validação opcional para PATCH)
	if validationErrors := validator.Validate(&request); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(w, validationErrors)
		return
	}

	// Converter DTO para map
	updates := request.ToMap()

	ctx := r.Context()
	updated, err := h.service.Patch(ctx, id, updates)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	// Converter model para DTO de resposta
	response := dto.FromModel(updated)

	utils.SuccessResponse(w, http.StatusOK, response)
}

// DeleteProduto deleta um produto
// DELETE /api/produtos/{id}
func (h *ProdutoHandler) DeleteProduto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, errors.ErrInvalidID)
		return
	}

	ctx := r.Context()
	if err := h.service.Delete(ctx, id); err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HealthCheckHandler gerencia o health check da API
type HealthCheckHandler struct {
	healthCheckFunc func(ctx context.Context) error
}

// NewHealthCheckHandler cria uma nova instância do HealthCheckHandler
func NewHealthCheckHandler(healthCheckFunc func(ctx context.Context) error) *HealthCheckHandler {
	return &HealthCheckHandler{
		healthCheckFunc: healthCheckFunc,
	}
}

// HealthCheck verifica o status da API e da conexão com o banco de dados
// GET /health
func (h *HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Verificar conexão com banco de dados se função disponível
	if h.healthCheckFunc != nil {
		if err := h.healthCheckFunc(ctx); err != nil {
			utils.JSONResponse(w, http.StatusServiceUnavailable, map[string]interface{}{
				"status":  "unhealthy",
				"message": "Conexão com banco de dados falhou",
				"error":   err.Error(),
			})
			return
		}
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"message": "API e banco de dados estão funcionando",
	})
}
