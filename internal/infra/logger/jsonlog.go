package logger

import (
	"encoding/json"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

type Logger interface {
	PrintInfo(message string, properties map[string]string)
	PrintError(err error, properties map[string]string)
	PrintFatal(err error, properties map[string]string)
}

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *logger {
	return &logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
}

func (l *logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level > LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ":unable to marshal log message:" + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

func (l *logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
