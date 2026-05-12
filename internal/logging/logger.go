package logging

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"
)

type Logger struct {
	nodeID string
	out    io.Writer
	mu     sync.Mutex
}

func New(nodeID string) *Logger {
	return &Logger{nodeID: nodeID, out: os.Stdout}
}

func NewWithWriter(nodeID string, out io.Writer) *Logger {
	return &Logger{nodeID: nodeID, out: out}
}

func (l *Logger) Info(message string, fields map[string]any) {
	l.write("info", message, fields)
}

func (l *Logger) Warn(message string, fields map[string]any) {
	l.write("warn", message, fields)
}

func (l *Logger) Error(message string, fields map[string]any) {
	l.write("error", message, fields)
}

func (l *Logger) write(level, message string, fields map[string]any) {
	event := map[string]any{
		"ts":      time.Now().UTC().Format(time.RFC3339Nano),
		"level":   level,
		"node_id": l.nodeID,
		"message": message,
	}
	for key, value := range fields {
		event[key] = value
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	_ = json.NewEncoder(l.out).Encode(event)
}
