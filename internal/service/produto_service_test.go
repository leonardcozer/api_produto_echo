package service

import (
	"context"
	"errors"
	"testing"

	apiErrors "api-go-arquitetura/internal/errors"
	"api-go-arquitetura/internal/model"
	"api-go-arquitetura/internal/repository"
)

// MockRepository é um mock do ProdutoRepository para testes
type MockRepository struct {
	produtos []model.Produto
	nextID   int
}

func NewMockRepository() repository.ProdutoRepository {
	return &MockRepository{
		produtos: make([]model.Produto, 0),
		nextID:   1,
	}
}

func (m *MockRepository) Create(ctx context.Context, produto model.Produto) (model.Produto, error) {
	produto.ID = m.nextID
	m.nextID++
	m.produtos = append(m.produtos, produto)
	return produto, nil
}

func (m *MockRepository) FindAll(ctx context.Context) ([]model.Produto, error) {
	return m.produtos, nil
}

func (m *MockRepository) FindByID(ctx context.Context, id int) (model.Produto, error) {
	for _, p := range m.produtos {
		if p.ID == id {
			return p, nil
		}
	}
	return model.Produto{}, errors.New("not found")
}

func (m *MockRepository) Update(ctx context.Context, id int, produto model.Produto) (model.Produto, error) {
	for i, p := range m.produtos {
		if p.ID == id {
			produto.ID = id
			m.produtos[i] = produto
			return produto, nil
		}
	}
	return model.Produto{}, errors.New("not found")
}

func (m *MockRepository) Patch(ctx context.Context, id int, updates map[string]interface{}) (model.Produto, error) {
	for i, p := range m.produtos {
		if p.ID == id {
			if nome, ok := updates["nome"].(string); ok {
				p.Nome = nome
			}
			if preco, ok := updates["preco"].(float64); ok {
				p.Preco = preco
			}
			if descricao, ok := updates["descricao"].(string); ok {
				p.Descricao = descricao
			}
			m.produtos[i] = p
			return p, nil
		}
	}
	return model.Produto{}, errors.New("not found")
}

func (m *MockRepository) Delete(ctx context.Context, id int) error {
	for i, p := range m.produtos {
		if p.ID == id {
			m.produtos = append(m.produtos[:i], m.produtos[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func TestProdutoService_Create(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockRepository()
	service := NewProdutoService(mockRepo)

	t.Run("deve criar produto com dados válidos", func(t *testing.T) {
		produto := model.Produto{
			Nome:      "Notebook",
			Preco:     3500.00,
			Descricao: "Notebook de alta performance",
		}

		result, err := service.Create(ctx, produto)
		if err != nil {
			t.Fatalf("Erro inesperado: %v", err)
		}

		if result.ID == 0 {
			t.Error("ID não foi atribuído")
		}
		if result.Nome != produto.Nome {
			t.Errorf("Nome esperado %s, obtido %s", produto.Nome, result.Nome)
		}
	})

	t.Run("deve retornar erro quando nome está vazio", func(t *testing.T) {
		produto := model.Produto{
			Preco:     3500.00,
			Descricao: "Produto sem nome",
		}

		_, err := service.Create(ctx, produto)
		if err == nil {
			t.Error("Esperado erro, mas nenhum erro foi retornado")
		}

		if !apiErrors.IsAPIError(err) {
			t.Error("Erro deveria ser do tipo APIError")
		}

		apiErr := apiErrors.AsAPIError(err)
		if apiErr.Code != "NOME_OBRIGATORIO" {
			t.Errorf("Código de erro esperado NOME_OBRIGATORIO, obtido %s", apiErr.Code)
		}
	})

	t.Run("deve retornar erro quando preço é inválido", func(t *testing.T) {
		produto := model.Produto{
			Nome:      "Notebook",
			Preco:     -100,
			Descricao: "Produto com preço negativo",
		}

		_, err := service.Create(ctx, produto)
		if err == nil {
			t.Error("Esperado erro, mas nenhum erro foi retornado")
		}

		apiErr := apiErrors.AsAPIError(err)
		if apiErr == nil || apiErr.Code != "PRECO_INVALIDO" {
			t.Errorf("Código de erro esperado PRECO_INVALIDO, obtido %v", apiErr)
		}
	})
}

func TestProdutoService_FindByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockRepository()
	service := NewProdutoService(mockRepo)

	// Criar produto de teste
	produto := model.Produto{
		Nome:      "Notebook",
		Preco:     3500.00,
		Descricao: "Notebook de alta performance",
	}
	created, _ := service.Create(ctx, produto)

	t.Run("deve encontrar produto existente", func(t *testing.T) {
		result, err := service.FindByID(ctx, created.ID)
		if err != nil {
			t.Fatalf("Erro inesperado: %v", err)
		}

		if result.ID != created.ID {
			t.Errorf("ID esperado %d, obtido %d", created.ID, result.ID)
		}
	})

	t.Run("deve retornar erro quando produto não existe", func(t *testing.T) {
		_, err := service.FindByID(ctx, 999)
		if err == nil {
			t.Error("Esperado erro, mas nenhum erro foi retornado")
		}

		apiErr := apiErrors.AsAPIError(err)
		if apiErr == nil || apiErr.Code != "PRODUTO_NOT_FOUND" {
			t.Errorf("Código de erro esperado PRODUTO_NOT_FOUND, obtido %v", apiErr)
		}
	})

	t.Run("deve retornar erro quando ID é inválido", func(t *testing.T) {
		_, err := service.FindByID(ctx, 0)
		if err == nil {
			t.Error("Esperado erro, mas nenhum erro foi retornado")
		}

		apiErr := apiErrors.AsAPIError(err)
		if apiErr == nil || apiErr.Code != "INVALID_ID" {
			t.Errorf("Código de erro esperado INVALID_ID, obtido %v", apiErr)
		}
	})
}

func TestProdutoService_Delete(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockRepository()
	service := NewProdutoService(mockRepo)

	// Criar produto de teste
	produto := model.Produto{
		Nome:      "Notebook",
		Preco:     3500.00,
		Descricao: "Notebook de alta performance",
	}
	created, _ := service.Create(ctx, produto)

	t.Run("deve deletar produto existente", func(t *testing.T) {
		err := service.Delete(ctx, created.ID)
		if err != nil {
			t.Fatalf("Erro inesperado: %v", err)
		}

		// Verificar se produto foi deletado
		_, err = service.FindByID(ctx, created.ID)
		if err == nil {
			t.Error("Produto deveria ter sido deletado")
		}
	})

	t.Run("deve retornar erro quando produto não existe", func(t *testing.T) {
		err := service.Delete(ctx, 999)
		if err == nil {
			t.Error("Esperado erro, mas nenhum erro foi retornado")
		}

		apiErr := apiErrors.AsAPIError(err)
		if apiErr == nil || apiErr.Code != "PRODUTO_NOT_FOUND" {
			t.Errorf("Código de erro esperado PRODUTO_NOT_FOUND, obtido %v", apiErr)
		}
	})
}

