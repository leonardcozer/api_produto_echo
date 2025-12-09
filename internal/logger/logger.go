package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// Log é a instância global do logger estruturado
	Log *logrus.Logger
	// lokiHook é o hook para enviar logs ao Loki
	lokiHook *LokiHook
)

func init() {
	Log = logrus.New()
	
	// Configurar formato JSON para produção (Loki requer JSON)
	if os.Getenv("LOG_FORMAT") == "json" || os.Getenv("LOKI_URL") != "" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	} else {
		// Formato texto para desenvolvimento
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
	
	// Configurar nível de log
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}
	
	// Output para stdout (sempre manter para logs locais)
	Log.SetOutput(os.Stdout)
	
	// Configurar hook do Loki se URL estiver configurada
	lokiURL := os.Getenv("LOKI_URL")
	lokiJob := os.Getenv("LOKI_JOB")
	if lokiJob == "" {
		lokiJob = "ARQUITETURA" // Valor padrão
	}
	
	if lokiURL != "" {
		lokiHook = NewLokiHook(lokiURL, lokiJob)
		if lokiHook != nil {
			Log.AddHook(lokiHook)
			// Forçar formato JSON quando Loki está ativo
			Log.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			})
		}
	}
}

// Shutdown encerra o logger e faz flush final dos logs para Loki
func Shutdown() {
	if lokiHook != nil {
		lokiHook.Stop()
	}
}

// WithField adiciona um campo ao logger
func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

// WithFields adiciona múltiplos campos ao logger
func WithFields(fields map[string]interface{}) *logrus.Entry {
	return Log.WithFields(logrus.Fields(fields))
}

// Debug registra mensagem de debug
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

// Debugf registra mensagem de debug formatada
func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

// Info registra mensagem de informação
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Infof registra mensagem de informação formatada
func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

// Warn registra mensagem de aviso
func Warn(args ...interface{}) {
	Log.Warn(args...)
}

// Warnf registra mensagem de aviso formatada
func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

// Error registra mensagem de erro
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Errorf registra mensagem de erro formatada
func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

// Fatal registra mensagem fatal e encerra o programa
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

// Fatalf registra mensagem fatal formatada e encerra o programa
func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

// Panic registra mensagem de panic
func Panic(args ...interface{}) {
	Log.Panic(args...)
}

// Panicf registra mensagem de panic formatada
func Panicf(format string, args ...interface{}) {
	Log.Panicf(format, args...)
}

