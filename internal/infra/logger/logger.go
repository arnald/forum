package logger

import (
	"fmt"
	"io"
	"os"
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
	case LevelOff:
		return "OFF"
	default:
		return ""
	}
}

type logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) Logger {
	return &logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *logger) PrintInfo(message string, properties map[string]string) {
	_, err := l.print(LevelInfo, message, properties)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write info log: %v\n", err)
	}
}

func (l *logger) PrintError(err error, properties map[string]string) {
	_, printErr := l.print(LevelError, err.Error(), properties)
	if printErr != nil {
		fmt.Fprintf(os.Stderr, "failed to write error log: %v\n", printErr)
	}
}

func (l *logger) PrintFatal(err error, properties map[string]string) {
	_, printErr := l.print(LevelFatal, err.Error(), properties)
	if printErr != nil {
		fmt.Fprintf(os.Stderr, "failed to write fatal log: %v\n", printErr)
	}
	os.Exit(1)
}

func (l *logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}

func (l *logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	timestamp := time.Now().UTC().Format(time.TimeOnly)

	logMsg := fmt.Sprintf("%-6s - [%s] - %s", timestamp, level.String(), message)

	if len(properties) > 0 {
		logMsg += " - "
		for key, value := range properties {
			logMsg += fmt.Sprintf("%s: %s; ", key, value)
		}
	}

	if level > LevelError {
		logMsg += "\nStack trace:\n" + string(debug.Stack())
	}

	logMsg += "\n"

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write([]byte(logMsg))
}
