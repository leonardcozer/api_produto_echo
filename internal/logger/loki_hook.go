package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// LokiHook é um hook do logrus que envia logs para Loki
type LokiHook struct {
	client      *http.Client
	url         string
	job         string
	batch       []lokiEntry
	batchMutex  sync.Mutex
	batchSize   int
	batchWait   time.Duration
	lastFlush   time.Time
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

// lokiEntry representa uma entrada de log para Loki
type lokiEntry struct {
	Timestamp string            `json:"ts"`
	Line      string            `json:"line"`
	Labels    map[string]string `json:"labels"`
}

// lokiStream representa um stream de logs para Loki
type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// lokiPayload representa o payload completo para Loki
type lokiPayload struct {
	Streams []lokiStream `json:"streams"`
}

// NewLokiHook cria um novo hook para Loki
func NewLokiHook(url, job string) *LokiHook {
	if url == "" {
		return nil
	}

	hook := &LokiHook{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		url:       url,
		job:       job,
		batch:     make([]lokiEntry, 0),
		batchSize: 10,              // Enviar em lotes de 10
		batchWait: 5 * time.Second, // Ou a cada 5 segundos
		lastFlush: time.Now(),
		stopChan:  make(chan struct{}),
	}

	// Iniciar goroutine para flush periódico
	hook.wg.Add(1)
	go hook.periodicFlush()

	return hook
}

// Levels retorna os níveis de log que o hook processa
func (h *LokiHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

// Fire é chamado quando um log é gerado
func (h *LokiHook) Fire(entry *logrus.Entry) error {
	if h == nil {
		return nil
	}

	// Converter entry para JSON
	var logLine bytes.Buffer
	encoder := json.NewEncoder(&logLine)
	encoder.SetEscapeHTML(false)
	
	logData := make(map[string]interface{})
	logData["level"] = entry.Level.String()
	logData["message"] = entry.Message
	logData["time"] = entry.Time.Format(time.RFC3339Nano)
	
	// Adicionar campos adicionais
	for k, v := range entry.Data {
		logData[k] = v
	}
	
	if err := encoder.Encode(logData); err != nil {
		return err
	}
	
	// Remover quebra de linha final do JSON
	logLineStr := logLine.String()
	if len(logLineStr) > 0 && logLineStr[len(logLineStr)-1] == '\n' {
		logLineStr = logLineStr[:len(logLineStr)-1]
	}

	// Criar labels para Loki
	labels := map[string]string{
		"job":      h.job,
		"level":    entry.Level.String(),
		"instance": getHostname(),
	}

	// Adicionar labels de campos específicos se existirem
	if method, ok := entry.Data["method"].(string); ok {
		labels["method"] = method
	}
	if path, ok := entry.Data["path"].(string); ok {
		labels["path"] = path
	}
	if statusCode, ok := entry.Data["status_code"].(int); ok {
		labels["status_code"] = fmt.Sprintf("%d", statusCode)
	}

	// Criar entrada
	lokiEntry := lokiEntry{
		Timestamp: fmt.Sprintf("%d", entry.Time.UnixNano()),
		Line:      logLineStr,
		Labels:    labels,
	}

	// Adicionar ao batch
	h.batchMutex.Lock()
	h.batch = append(h.batch, lokiEntry)
	batchSize := len(h.batch)
	h.batchMutex.Unlock()

	// Se o batch atingir o tamanho máximo, enviar imediatamente
	if batchSize >= h.batchSize {
		go h.flush()
	}

	return nil
}

// flush envia o batch atual para Loki
func (h *LokiHook) flush() error {
	h.batchMutex.Lock()
	if len(h.batch) == 0 {
		h.batchMutex.Unlock()
		return nil
	}

	// Copiar batch e limpar
	batch := make([]lokiEntry, len(h.batch))
	copy(batch, h.batch)
	h.batch = h.batch[:0]
	h.batchMutex.Unlock()

	// Agrupar por labels
	streamsMap := make(map[string]*lokiStream)
	for _, entry := range batch {
		// Criar chave única para labels
		labelKey := fmt.Sprintf("%v", entry.Labels)
		
		stream, exists := streamsMap[labelKey]
		if !exists {
			stream = &lokiStream{
				Stream: entry.Labels,
				Values: make([][]string, 0),
			}
			streamsMap[labelKey] = stream
		}
		
		stream.Values = append(stream.Values, []string{
			entry.Timestamp,
			entry.Line,
		})
	}

	// Converter para array
	streams := make([]lokiStream, 0, len(streamsMap))
	for _, stream := range streamsMap {
		streams = append(streams, *stream)
	}

	payload := lokiPayload{
		Streams: streams,
	}

	// Serializar payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar payload: %w", err)
	}

	// Enviar para Loki
	req, err := http.NewRequest("POST", h.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar para Loki: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Loki retornou status %d", resp.StatusCode)
	}

	h.lastFlush = time.Now()
	return nil
}

// periodicFlush envia logs periodicamente
func (h *LokiHook) periodicFlush() {
	defer h.wg.Done()

	ticker := time.NewTicker(h.batchWait)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.flush()
		case <-h.stopChan:
			// Flush final antes de sair
			h.flush()
			return
		}
	}
}

// Stop para o hook e faz flush final
func (h *LokiHook) Stop() {
	if h == nil {
		return
	}
	close(h.stopChan)
	h.wg.Wait()
}

// getHostname retorna o hostname da máquina
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

