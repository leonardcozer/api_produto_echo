package database

import (
	"context"
	"fmt"
	"time"

	"api-go-arquitetura/internal/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

// RetryOptions configura opções de retry
type RetryOptions struct {
	MaxAttempts int           // Número máximo de tentativas
	InitialDelay time.Duration // Delay inicial entre tentativas
	MaxDelay     time.Duration // Delay máximo entre tentativas
	Multiplier   float64      // Multiplicador para backoff exponencial
}

// DefaultRetryOptions retorna opções padrão de retry
func DefaultRetryOptions() RetryOptions {
	return RetryOptions{
		MaxAttempts: 3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
	}
}

// RetryableError verifica se um erro é retryable
func RetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Erros de rede são retryable
	if mongo.IsNetworkError(err) {
		return true
	}

	// Timeout errors são retryable
	if mongo.IsTimeout(err) {
		return true
	}

	// Erros de conexão são retryable
	if err == mongo.ErrClientDisconnected {
		return true
	}

	return false
}

// Retry executa uma função com retry logic
func Retry(ctx context.Context, fn func() error, opts RetryOptions) error {
	var lastErr error
	delay := opts.InitialDelay

	for attempt := 1; attempt <= opts.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Se não é um erro retryable, retornar imediatamente
		if !RetryableError(err) {
			return err
		}

		// Se é a última tentativa, não esperar
		if attempt == opts.MaxAttempts {
			break
		}

		// Log da tentativa
		logger.WithFields(map[string]interface{}{
			"attempt": attempt,
			"max_attempts": opts.MaxAttempts,
			"error": err.Error(),
			"delay_ms": delay.Milliseconds(),
		}).Warn("Retry attempt failed, retrying...")

		// Aguardar antes da próxima tentativa
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		// Calcular próximo delay (backoff exponencial)
		delay = time.Duration(float64(delay) * opts.Multiplier)
		if delay > opts.MaxDelay {
			delay = opts.MaxDelay
		}
	}

	return fmt.Errorf("max retry attempts (%d) reached: %w", opts.MaxAttempts, lastErr)
}

// RetryWithResult executa uma função com retry logic que retorna um resultado
func RetryWithResult[T any](ctx context.Context, fn func() (T, error), opts RetryOptions) (T, error) {
	var zero T
	var lastErr error
	delay := opts.InitialDelay

	for attempt := 1; attempt <= opts.MaxAttempts; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}

		lastErr = err

		// Se não é um erro retryable, retornar imediatamente
		if !RetryableError(err) {
			return zero, err
		}

		// Se é a última tentativa, não esperar
		if attempt == opts.MaxAttempts {
			break
		}

		// Log da tentativa
		logger.WithFields(map[string]interface{}{
			"attempt": attempt,
			"max_attempts": opts.MaxAttempts,
			"error": err.Error(),
			"delay_ms": delay.Milliseconds(),
		}).Warn("Retry attempt failed, retrying...")

		// Aguardar antes da próxima tentativa
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		case <-time.After(delay):
		}

		// Calcular próximo delay (backoff exponencial)
		delay = time.Duration(float64(delay) * opts.Multiplier)
		if delay > opts.MaxDelay {
			delay = opts.MaxDelay
		}
	}

	return zero, fmt.Errorf("max retry attempts (%d) reached: %w", opts.MaxAttempts, lastErr)
}

