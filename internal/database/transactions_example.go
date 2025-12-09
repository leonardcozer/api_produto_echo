package database

// Este arquivo contém exemplos de uso de transações
// Não é compilado, apenas para referência

/*
Exemplo de uso de transações:

import (
	"context"
	"api-go-arquitetura/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func exemploTransacao(ctx context.Context, client *mongo.Client) error {
	// Iniciar transação
	tx, cancel, err := database.StartTransaction(ctx, client)
	if err != nil {
		return err
	}
	defer cancel() // Importante: sempre cancelar o contexto
	defer tx.End() // Importante: sempre finalizar a sessão

	// Executar operações dentro da transação
	err = tx.WithTransaction(func(sc mongo.SessionContext) error {
		// Todas as operações de banco devem usar sc como contexto
		// Exemplo: inserir em múltiplas coleções
		
		// Operação 1
		// _, err := collection1.InsertOne(sc, document1)
		// if err != nil {
		//     return err // Isso fará rollback automático
		// }
		
		// Operação 2
		// _, err := collection2.InsertOne(sc, document2)
		// if err != nil {
		//     return err // Isso fará rollback automático
		// }
		
		return nil // Sucesso - commit automático
	})

	if err != nil {
		// Transação foi abortada automaticamente
		return err
	}

	// Transação foi commitada com sucesso
	return nil
}
*/

