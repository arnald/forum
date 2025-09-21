/*
Package logger provides structured logging functionality for the forum application.

This package implements a custom logger with:
- Multiple log levels (INFO, ERROR, FATAL, OFF)
- Structured logging with key-value properties
- Thread-safe logging operations
- Configurable output destinations
- Stack trace capture for fatal errors
- RFC3339 timestamp formatting

The logger is designed to be simple, efficient, and suitable for production use
while providing enough flexibility for development and debugging.
*/
package logger

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Level represents the severity level of log messages
// Higher levels indicate more severe events that require attention
type Level int8

// Log level constants defining the severity hierarchy
// Logs are only printed if they meet or exceed the configured minimum level
const (
	LevelInfo  Level = iota // Informational messages for normal operation
	LevelError              // Error conditions that don't stop the application
	LevelFatal              // Critical errors that cause application termination
	LevelOff                // Disable all logging output
)

// Logger interface defines the contract for logging operations
// This interface allows for easy testing and different logger implementations
type Logger interface {
	PrintInfo(message string, properties map[string]string)  // Log informational messages
	PrintError(err error, properties map[string]string)     // Log error conditions
	PrintFatal(err error, properties map[string]string)     // Log fatal errors and exit
}

// String returns the string representation of a log level
// This is used in log message formatting to show the severity level
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"   // Informational level
	case LevelError:
		return "ERROR"  // Error level
	case LevelFatal:
		return "FATAL"  // Fatal level
	case LevelOff:
		return "OFF"    // Logging disabled
	default:
		return ""       // Unknown level
	}
}

// logger is the concrete implementation of the Logger interface
// It provides thread-safe logging operations with configurable output and levels
type logger struct {
	out      io.Writer   // Output destination (stdout, file, etc.)
	minLevel Level       // Minimum level to log (filters out lower levels)
	mu       sync.Mutex  // Mutex for thread-safe writing
}

// New creates a new logger instance with specified output and minimum level
// This factory function configures the logger for the application's needs
//
// Parameters:
//   - out: io.Writer where log messages will be written (os.Stdout, file, etc.)
//   - minLevel: Minimum log level to output (filters lower-level messages)
//
// Returns:
//   - Logger: Configured logger ready for use
//
// Usage:
//   logger := New(os.Stdout, LevelInfo)  // Log info and above to stdout
//   logger := New(file, LevelError)      // Log only errors and fatal to file
func New(out io.Writer, minLevel Level) Logger {
	return &logger{
		out:      out,      // Set output destination
		minLevel: minLevel, // Set filtering level
	}
}

// PrintInfo logs informational messages with optional structured properties
// These messages indicate normal application operation and flow
//
// Parameters:
//   - message: The main log message describing the event
//   - properties: Optional key-value pairs providing additional context
//
// Info logs are useful for:
// - Application startup/shutdown events
// - User actions and business operations
// - Performance metrics and statistics
func (l *logger) PrintInfo(message string, properties map[string]string) {
	_, err := l.print(LevelInfo, message, properties)
	if err != nil {
		// Fallback to stderr if regular logging fails
		fmt.Fprintf(os.Stderr, "failed to write info log: %v\n", err)
	}
}

// PrintError logs error conditions with optional structured properties
// These messages indicate problems that occurred but didn't stop the application
//
// Parameters:
//   - err: The error that occurred (provides error message)
//   - properties: Optional key-value pairs providing additional context
//
// Error logs are useful for:
// - Failed operations that were handled gracefully
// - Invalid user input or requests
// - External service failures
// - Resource constraints or temporary issues
func (l *logger) PrintError(err error, properties map[string]string) {
	_, printErr := l.print(LevelError, err.Error(), properties)
	if printErr != nil {
		// Fallback to stderr if regular logging fails
		fmt.Fprintf(os.Stderr, "failed to write error log: %v\n", printErr)
	}
}

// PrintFatal logs critical errors and terminates the application
// These messages indicate severe problems that prevent continued operation
//
// Parameters:
//   - err: The fatal error that occurred
//   - properties: Optional key-value pairs providing additional context
//
// Fatal logs include:
// - Full stack trace for debugging
// - Immediate application termination (os.Exit(1))
//
// Use for:
// - Database connection failures
// - Critical configuration errors
// - Security violations
// - Unrecoverable system errors
func (l *logger) PrintFatal(err error, properties map[string]string) {
	_, printErr := l.print(LevelFatal, err.Error(), properties)
	if printErr != nil {
		// Fallback to stderr if regular logging fails
		fmt.Fprintf(os.Stderr, "failed to write fatal log: %v\n", printErr)
	}
	os.Exit(1) // Terminate application immediately
}

// Write implements io.Writer interface for integration with other tools
// This allows the logger to be used as a writer for HTTP servers, etc.
//
// Parameters:
//   - message: Raw message bytes to log
//
// Returns:
//   - n: Number of bytes written
//   - err: Any error that occurred during writing
//
// Messages written through this interface are logged at ERROR level
func (l *logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}

// print is the core logging function that formats and writes log messages
// This method handles the actual message formatting, filtering, and output
//
// Parameters:
//   - level: Severity level of the message
//   - message: Main message content to log
//   - properties: Optional structured data to include
//
// Returns:
//   - int: Number of bytes written
//   - error: Any error that occurred during writing
//
// Message format: "TIMESTAMP - [LEVEL] - MESSAGE - key: value; key: value;"
// Fatal messages also include full stack traces for debugging
func (l *logger) print(level Level, message string, properties map[string]string) (int, error) {
	// Filter out messages below minimum level
	if level < l.minLevel {
		return 0, nil
	}

	// Create RFC3339 timestamp in UTC for consistency
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Format base log message with timestamp and level
	logMsg := fmt.Sprintf("%-6s - [%s] - %s", timestamp, level.String(), message)

	// Add structured properties if provided
	if len(properties) > 0 {
		logMsg += " - "
		for key, value := range properties {
			logMsg += fmt.Sprintf("%s: %s; ", key, value)
		}
	}

	// Add stack trace for fatal errors to aid debugging
	if level > LevelError {
		logMsg += "\nStack trace:\n" + string(debug.Stack())
	}

	logMsg += "\n" // Ensure each log message ends with newline

	// Use mutex to ensure thread-safe writing
	l.mu.Lock()
	defer l.mu.Unlock()

	// Write formatted message to configured output
	return l.out.Write([]byte(logMsg))
}
