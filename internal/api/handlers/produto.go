package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

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
func (h *ProdutoHandler) GetProdutos(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Parse de parâmetros de paginação
	pagination := dto.PaginationRequest{
		Page:     getIntQueryEcho(c, "page", 1),
		PageSize: getIntQueryEcho(c, "pageSize", 10),
	}

	// Parse de filtros
	filter := dto.FilterRequest{
		Nome:      getStringQueryEcho(c, "nome"),
		PrecoMin:  getFloatQueryEcho(c, "precoMin"),
		PrecoMax:  getFloatQueryEcho(c, "precoMax"),
		Descricao: getStringQueryEcho(c, "descricao"),
	}

	// Parse de ordenação
	sort := dto.GetSortFromQuery(
		c.QueryParam("sort"),
		c.QueryParam("order"),
	)
	
	// Validar ordenação
	if validationErrors := validator.Validate(&sort); len(validationErrors) > 0 {
		return utils.EchoValidationErrorResponse(c, validationErrors)
	}

	// Se não há filtros e paginação padrão, usar método antigo para compatibilidade
	if filter.IsEmpty() && pagination.Page == 1 && pagination.PageSize == 10 && sort.Field == "" {
		// Verificar se há parâmetros de query explícitos
		if c.QueryParam("page") == "" && c.QueryParam("pageSize") == "" && c.QueryParam("sort") == "" {
			// Usar método antigo (sem paginação)
			produtos, err := h.service.FindAll(ctx)
			if err != nil {
				return utils.EchoErrorResponse(c, errors.WrapError(err, errors.ErrDatabase))
			}
			response := dto.ToProdutoListResponse(produtos)
			return utils.EchoSuccessResponse(c, http.StatusOK, response)
		}
	}

	// Usar método paginado
	produtos, paginationResp, err := h.service.FindAllPaginated(ctx, pagination, filter, sort)
	if err != nil {
		return utils.EchoErrorResponse(c, errors.WrapError(err, errors.ErrDatabase))
	}

	// Converter models para DTOs
	produtosDTO := dto.FromModelList(produtos)
	response := dto.ToPaginatedResponse(produtosDTO, paginationResp)

	return utils.EchoSuccessResponse(c, http.StatusOK, response)
}

// getIntQueryEcho obtém um parâmetro de query como int usando Echo
func getIntQueryEcho(c echo.Context, key string, defaultValue int) int {
	value := c.QueryParam(key)
	if value == "" {
		return defaultValue
	}
	var result int
	if _, err := fmt.Sscanf(value, "%d", &result); err != nil {
		return defaultValue
	}
	return result
}

// getStringQueryEcho obtém um parâmetro de query como string usando Echo (retorna nil se vazio)
func getStringQueryEcho(c echo.Context, key string) *string {
	value := c.QueryParam(key)
	if value == "" {
		return nil
	}
	return &value
}

// getFloatQueryEcho obtém um parâmetro de query como float64 usando Echo (retorna nil se vazio)
func getFloatQueryEcho(c echo.Context, key string) *float64 {
	value := c.QueryParam(key)
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
func (h *ProdutoHandler) GetProduto(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.EchoErrorResponse(c, errors.ErrInvalidID)
	}

	ctx := c.Request().Context()
	produto, err := h.service.FindByID(ctx, id)
	if err != nil {
		if errors.IsAPIError(err) {
			return utils.EchoErrorResponse(c, err)
		}
		return utils.EchoErrorResponse(c, errors.ErrProdutoNotFound)
	}

	// Converter model para DTO
	response := dto.FromModel(produto)

	return utils.EchoSuccessResponse(c, http.StatusOK, response)
}

// CreateProduto cria um novo produto
// POST /api/produtos
func (h *ProdutoHandler) CreateProduto(c echo.Context) error {
	var request dto.CreateProdutoRequest
	
	// Decodificar JSON usando Echo
	if err := c.Bind(&request); err != nil {
		return utils.EchoBadRequestResponse(c, "Erro ao decodificar JSON: "+err.Error())
	}

	// Validar DTO
	if validationErrors := validator.Validate(&request); len(validationErrors) > 0 {
		return utils.EchoValidationErrorResponse(c, validationErrors)
	}

	// Converter DTO para model
	produto := request.ToModel()

	ctx := c.Request().Context()
	created, err := h.service.Create(ctx, produto)
	if err != nil {
		return utils.EchoErrorResponse(c, err)
	}

	// Converter model para DTO de resposta
	response := dto.FromModel(created)

	return utils.EchoSuccessResponse(c, http.StatusCreated, response)
}

// UpdateProduto atualiza um produto completamente
// PUT /api/produtos/{id}
func (h *ProdutoHandler) UpdateProduto(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.EchoErrorResponse(c, errors.ErrInvalidID)
	}

	var request dto.UpdateProdutoRequest
	
	// Decodificar JSON usando Echo
	if err := c.Bind(&request); err != nil {
		return utils.EchoBadRequestResponse(c, "Erro ao decodificar JSON: "+err.Error())
	}

	// Validar DTO
	if validationErrors := validator.Validate(&request); len(validationErrors) > 0 {
		return utils.EchoValidationErrorResponse(c, validationErrors)
	}

	// Converter DTO para model
	produto := request.ToModel()

	ctx := c.Request().Context()
	updated, err := h.service.Update(ctx, id, produto)
	if err != nil {
		return utils.EchoErrorResponse(c, err)
	}

	// Converter model para DTO de resposta
	response := dto.FromModel(updated)

	return utils.EchoSuccessResponse(c, http.StatusOK, response)
}

// PatchProduto atualiza um produto parcialmente
// PATCH /api/produtos/{id}
func (h *ProdutoHandler) PatchProduto(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.EchoErrorResponse(c, errors.ErrInvalidID)
	}

	var request dto.PatchProdutoRequest
	
	// Decodificar JSON usando Echo
	if err := c.Bind(&request); err != nil {
		return utils.EchoBadRequestResponse(c, "Erro ao decodificar JSON: "+err.Error())
	}

	// Validar DTO (validação opcional para PATCH)
	if validationErrors := validator.Validate(&request); len(validationErrors) > 0 {
		return utils.EchoValidationErrorResponse(c, validationErrors)
	}

	// Converter DTO para map
	updates := request.ToMap()

	ctx := c.Request().Context()
	updated, err := h.service.Patch(ctx, id, updates)
	if err != nil {
		return utils.EchoErrorResponse(c, err)
	}

	// Converter model para DTO de resposta
	response := dto.FromModel(updated)

	return utils.EchoSuccessResponse(c, http.StatusOK, response)
}

// DeleteProduto deleta um produto
// DELETE /api/produtos/{id}
func (h *ProdutoHandler) DeleteProduto(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.EchoErrorResponse(c, errors.ErrInvalidID)
	}

	ctx := c.Request().Context()
	if err := h.service.Delete(ctx, id); err != nil {
		return utils.EchoErrorResponse(c, err)
	}

	return c.NoContent(http.StatusNoContent)
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
func (h *HealthCheckHandler) HealthCheck(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Verificar conexão com banco de dados se função disponível
	if h.healthCheckFunc != nil {
		if err := h.healthCheckFunc(ctx); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
				"status":  "unhealthy",
				"message": "Conexão com banco de dados falhou",
				"error":   err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"message": "API e banco de dados estão funcionando",
	})
}
