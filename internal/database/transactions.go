package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// TransactionSession encapsula uma sessão de transação do MongoDB
type TransactionSession struct {
	session mongo.Session
	ctx     context.Context
}

// StartTransaction inicia uma nova transação
func StartTransaction(ctx context.Context, client *mongo.Client) (*TransactionSession, func(), error) {
	session, err := client.StartSession()
	if err != nil {
		return nil, nil, err
	}

	// Criar contexto com timeout para a transação
	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

	return &TransactionSession{
		session: session,
		ctx:     txCtx,
	}, cancel, nil
}

// WithTransaction executa uma função dentro de uma transação
func (ts *TransactionSession) WithTransaction(fn func(mongo.SessionContext) error) error {
	return mongo.WithSession(ts.ctx, ts.session, func(sc mongo.SessionContext) error {
		if err := ts.session.StartTransaction(); err != nil {
			return err
		}

		// Executar a função
		if err := fn(sc); err != nil {
			// Em caso de erro, abortar a transação
			if abortErr := ts.session.AbortTransaction(sc); abortErr != nil {
				return abortErr
			}
			return err
		}

		// Se tudo correu bem, commitar a transação
		return ts.session.CommitTransaction(sc)
	})
}

// Abort aborta a transação atual
func (ts *TransactionSession) Abort() error {
	return ts.session.AbortTransaction(ts.ctx)
}

// End finaliza a sessão
func (ts *TransactionSession) End() {
	ts.session.EndSession(ts.ctx)
}

// GetContext retorna o contexto da sessão
func (ts *TransactionSession) GetContext() context.Context {
	return ts.ctx
}

// GetSession retorna a sessão do MongoDB
func (ts *TransactionSession) GetSession() mongo.Session {
	return ts.session
}

