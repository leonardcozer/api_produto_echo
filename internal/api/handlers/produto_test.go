package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"api-go-arquitetura/internal/dto"
	apiErrors "api-go-arquitetura/internal/errors"
	"api-go-arquitetura/internal/model"
	"api-go-arquitetura/internal/service"
)

// MockProdutoService é um mock do ProdutoService para testes
type MockProdutoService struct {
	produtos []model.Produto
	nextID   int
}

func NewMockProdutoService() *MockProdutoService {
	return &MockProdutoService{
		produtos: make([]model.Produto, 0),
		nextID:   1,
	}
}

func (m *MockProdutoService) Create(ctx context.Context, produto model.Produto) (model.Produto, error) {
	produto.ID = m.nextID
	m.nextID++
	m.produtos = append(m.produtos, produto)
	return produto, nil
}

func (m *MockProdutoService) FindAll(ctx context.Context) ([]model.Produto, error) {
	return m.produtos, nil
}

func (m *MockProdutoService) FindByID(ctx context.Context, id int) (model.Produto, error) {
	for _, p := range m.produtos {
		if p.ID == id {
			return p, nil
		}
	}
	return model.Produto{}, apiErrors.ErrProdutoNotFound
}

func (m *MockProdutoService) Update(ctx context.Context, id int, produto model.Produto) (model.Produto, error) {
	for i, p := range m.produtos {
		if p.ID == id {
			produto.ID = id
			m.produtos[i] = produto
			return produto, nil
		}
	}
	return model.Produto{}, apiErrors.ErrProdutoNotFound
}

func (m *MockProdutoService) Patch(ctx context.Context, id int, updates map[string]interface{}) (model.Produto, error) {
	for i, p := range m.produtos {
		if p.ID == id {
			if nome, ok := updates["nome"].(string); ok {
				p.Nome = nome
			}
			if preco, ok := updates["preco"].(float64); ok {
				p.Preco = preco
			}
			m.produtos[i] = p
			return p, nil
		}
	}
	return model.Produto{}, apiErrors.ErrProdutoNotFound
}

func (m *MockProdutoService) Delete(ctx context.Context, id int) error {
	for i, p := range m.produtos {
		if p.ID == id {
			m.produtos = append(m.produtos[:i], m.produtos[i+1:]...)
			return nil
		}
	}
	return apiErrors.ErrProdutoNotFound
}

func TestProdutoHandler_CreateProduto(t *testing.T) {
	mockService := NewMockProdutoService()
	handler := NewProdutoHandler(mockService)
	e := echo.New()

	t.Run("deve criar produto com dados válidos", func(t *testing.T) {
		requestBody := dto.CreateProdutoRequest{
			Nome:      "Notebook",
			Preco:     3500.00,
			Descricao: "Notebook de alta performance",
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/api/produtos", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateProduto(c)
		if err != nil {
			t.Fatalf("Erro inesperado: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Status esperado %d, obtido %d", http.StatusCreated, rec.Code)
		}

		var response dto.ProdutoResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("Erro ao decodificar resposta: %v", err)
		}

		if response.Nome != requestBody.Nome {
			t.Errorf("Nome esperado %s, obtido %s", requestBody.Nome, response.Nome)
		}
	})

	t.Run("deve retornar erro quando dados são inválidos", func(t *testing.T) {
		requestBody := dto.CreateProdutoRequest{
			Nome: "", // Nome vazio
			Preco: 3500.00,
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/api/produtos", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler.CreateProduto(c)

		// Pode ser 400 (Bad Request) ou 422 (Unprocessable Entity) dependendo da validação
		if rec.Code != http.StatusBadRequest && rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("Status esperado %d ou %d, obtido %d", http.StatusBadRequest, http.StatusUnprocessableEntity, rec.Code)
		}
	})
}

func TestProdutoHandler_GetProduto(t *testing.T) {
	mockService := NewMockProdutoService()
	handler := NewProdutoHandler(mockService)
	e := echo.New()

	// Criar produto de teste
	produto := model.Produto{
		Nome:      "Notebook",
		Preco:     3500.00,
		Descricao: "Notebook de alta performance",
	}
	created, _ := mockService.Create(context.Background(), produto)

	t.Run("deve retornar produto existente", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/produtos/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/produtos/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.GetProduto(c)
		if err != nil {
			t.Fatalf("Erro inesperado: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Status esperado %d, obtido %d", http.StatusOK, rec.Code)
		}

		var response dto.ProdutoResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("Erro ao decodificar resposta: %v", err)
		}

		if response.ID != created.ID {
			t.Errorf("ID esperado %d, obtido %d", created.ID, response.ID)
		}
	})

	t.Run("deve retornar erro quando produto não existe", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/produtos/999", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/produtos/:id")
		c.SetParamNames("id")
		c.SetParamValues("999")

		handler.GetProduto(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Status esperado %d, obtido %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("deve retornar erro quando ID é inválido", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/produtos/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/produtos/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		handler.GetProduto(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Status esperado %d, obtido %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestProdutoHandler_GetProdutos(t *testing.T) {
	mockService := NewMockProdutoService()
	handler := NewProdutoHandler(mockService)
	e := echo.New()

	// Criar produtos de teste
	mockService.Create(context.Background(), model.Produto{
		Nome:  "Notebook",
		Preco: 3500.00,
	})
	mockService.Create(context.Background(), model.Produto{
		Nome:  "Mouse",
		Preco: 50.00,
	})

	t.Run("deve retornar lista de produtos", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/produtos", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetProdutos(c)
		if err != nil {
			t.Fatalf("Erro inesperado: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Status esperado %d, obtido %d", http.StatusOK, rec.Code)
		}

		var response dto.ProdutoListResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("Erro ao decodificar resposta: %v", err)
		}

		if response.Total != 2 {
			t.Errorf("Total esperado 2, obtido %d", response.Total)
		}

		if len(response.Produtos) != 2 {
			t.Errorf("Quantidade de produtos esperada 2, obtida %d", len(response.Produtos))
		}
	})
}

