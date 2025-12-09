package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api-go-arquitetura/internal/logger"
	"api-go-arquitetura/internal/metrics"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// ErrConnectionFailed é retornado quando não é possível conectar ao MongoDB
	ErrConnectionFailed = errors.New("falha ao conectar ao MongoDB")
	// ErrPingFailed é retornado quando o ping ao MongoDB falha
	ErrPingFailed = errors.New("falha ao verificar conexão com MongoDB (ping)")
	// ErrInvalidURI é retornado quando a URI do MongoDB é inválida
	ErrInvalidURI = errors.New("URI do MongoDB inválida")
)

// ConnectOptions contém opções para conexão com o MongoDB
type ConnectOptions struct {
	URI            string
	ConnectTimeout time.Duration
	MaxPoolSize    uint64
	MinPoolSize    uint64
}

// DefaultConnectOptions retorna opções padrão para conexão
func DefaultConnectOptions(uri string) ConnectOptions {
	return ConnectOptions{
		URI:            uri,
		ConnectTimeout: 10 * time.Second,
		MaxPoolSize:    100,
		MinPoolSize:    10,
	}
}

// Connect estabelece uma conexão com o MongoDB e retorna o cliente
// Retorna erro se a conexão falhar ou se o ping não funcionar
func Connect(opts ConnectOptions) (*mongo.Client, error) {
	if opts.URI == "" {
		return nil, fmt.Errorf("%w: URI não pode ser vazia", ErrInvalidURI)
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.ConnectTimeout)
	defer cancel()

	// Configurar opções do cliente
	clientOptions := options.Client().
		ApplyURI(opts.URI).
		SetMaxPoolSize(opts.MaxPoolSize).
		SetMinPoolSize(opts.MinPoolSize).
		SetConnectTimeout(opts.ConnectTimeout).
		SetServerSelectionTimeout(5 * time.Second)

	// Tentar conectar
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	// Verificar se a conexão realmente funciona (ping)
	if err := Ping(ctx, client); err != nil {
		// Se o ping falhar, fechar a conexão e retornar erro
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("%w: %v", ErrPingFailed, err)
	}

	logger.WithField("uri", opts.URI).Info("Conexão com MongoDB estabelecida com sucesso")
	
	// Atualizar métricas de conexão
	updateConnectionMetrics(client, opts.MaxPoolSize)
	
	return client, nil
}

// updateConnectionMetrics atualiza as métricas de conexão do banco de dados
func updateConnectionMetrics(client *mongo.Client, maxPoolSize uint64) {
	// Obter estatísticas do pool de conexões
	stats := client.NumberSessionsInProgress()
	
	// Atualizar métricas
	metrics.SetDatabaseConnections("active", float64(stats))
	metrics.SetDatabaseConnections("total", float64(maxPoolSize))
	metrics.SetDatabaseConnections("idle", float64(maxPoolSize)-float64(stats))
}

// Ping verifica se a conexão com o MongoDB está funcionando
func Ping(ctx context.Context, client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("ping falhou: %w", err)
	}
	return nil
}

// Disconnect fecha a conexão com o MongoDB de forma adequada
func Disconnect(ctx context.Context, client *mongo.Client) error {
	if client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("erro ao desconectar do MongoDB: %w", err)
	}

	logger.Info("Conexão com MongoDB fechada com sucesso")
	return nil
}

// CreateIndexes cria índices otimizados para a coleção de produtos
func CreateIndexes(ctx context.Context, client *mongo.Client, database, collection string) error {
	col := client.Database(database).Collection(collection)

	// Índice único no campo ID
	idIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("idx_id"),
	}

	// Índice de texto para busca por nome (case-insensitive)
	nomeIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "nome", Value: 1}},
		Options: options.Index().SetName("idx_nome"),
	}

	// Índice para busca por descrição
	descricaoIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "descricao", Value: 1}},
		Options: options.Index().SetName("idx_descricao"),
	}

	// Índice composto para filtros de preço
	precoIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "preco", Value: 1}},
		Options: options.Index().SetName("idx_preco"),
	}

	// Criar todos os índices
	indexes := []mongo.IndexModel{idIndex, nomeIndex, descricaoIndex, precoIndex}
	_, err := col.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("erro ao criar índices: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"database":   database,
		"collection": collection,
		"indexes":    len(indexes),
	}).Info("Índices criados com sucesso")

	return nil
}

// HealthCheck verifica a saúde da conexão com o MongoDB
func HealthCheck(ctx context.Context, client *mongo.Client) error {
	if client == nil {
		return errors.New("cliente MongoDB não está inicializado")
	}

	// Verificar se o cliente ainda está conectado
	if err := Ping(ctx, client); err != nil {
		return fmt.Errorf("health check falhou: %w", err)
	}

	return nil
}

// GetDatabase retorna uma referência ao banco de dados com tratamento de erro
func GetDatabase(client *mongo.Client, dbName string) (*mongo.Database, error) {
	if client == nil {
		return nil, errors.New("cliente MongoDB não está inicializado")
	}

	if dbName == "" {
		return nil, errors.New("nome do banco de dados não pode ser vazio")
	}

	return client.Database(dbName), nil
}

// GetCollection retorna uma referência à coleção com tratamento de erro
func GetCollection(client *mongo.Client, dbName, collectionName string) (*mongo.Collection, error) {
	if client == nil {
		return nil, errors.New("cliente MongoDB não está inicializado")
	}

	if dbName == "" {
		return nil, errors.New("nome do banco de dados não pode ser vazio")
	}

	if collectionName == "" {
		return nil, errors.New("nome da coleção não pode ser vazio")
	}

	db, err := GetDatabase(client, dbName)
	if err != nil {
		return nil, err
	}

	return db.Collection(collectionName), nil
}