package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger interface defines logging operations
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	WithFields(fields ...Field) Logger
	WithError(err error) Logger
	WithContext(ctx string) Logger
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// NewField creates a new field
func NewField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	File      string                 `json:"file,omitempty"`
	Line      int                    `json:"line,omitempty"`
	Function  string                 `json:"function,omitempty"`
}

// StructuredLogger implements structured logging
type StructuredLogger struct {
	level    LogLevel
	output   io.Writer
	fields   map[string]interface{}
	context  string
	caller   bool
}

// NewStructuredLogger creates a new structured logger
func NewStructuredLogger(level LogLevel, output io.Writer) *StructuredLogger {
	if output == nil {
		output = os.Stdout
	}
	
	return &StructuredLogger{
		level:  level,
		output: output,
		fields: make(map[string]interface{}),
	}
}

// NewDefaultLogger creates a logger with default settings
func NewDefaultLogger() *StructuredLogger {
	return NewStructuredLogger(INFO, os.Stdout)
}

// NewDevelopmentLogger creates a logger for development
func NewDevelopmentLogger() *StructuredLogger {
	logger := NewStructuredLogger(DEBUG, os.Stdout)
	logger.caller = true
	return logger
}

// NewProductionLogger creates a logger for production
func NewProductionLogger() *StructuredLogger {
	return NewStructuredLogger(INFO, os.Stdout)
}

// SetLevel sets the logging level
func (l *StructuredLogger) SetLevel(level LogLevel) {
	l.level = level
}

// SetOutput sets the output writer
func (l *StructuredLogger) SetOutput(output io.Writer) {
	l.output = output
}

// EnableCaller enables caller information in logs
func (l *StructuredLogger) EnableCaller() {
	l.caller = true
}

// DisableCaller disables caller information in logs
func (l *StructuredLogger) DisableCaller() {
	l.caller = false
}

// log writes a log entry
func (l *StructuredLogger) log(level LogLevel, msg string, fields ...Field) {
	if level < l.level {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level.String(),
		Message:   msg,
		Fields:    make(map[string]interface{}),
	}

	// Add context if set
	if l.context != "" {
		entry.Fields["context"] = l.context
	}

	// Add existing fields
	for k, v := range l.fields {
		entry.Fields[k] = v
	}

	// Add new fields
	for _, field := range fields {
		entry.Fields[field.Key] = field.Value
	}

	// Add caller information if enabled
	if l.caller {
		if pc, file, line, ok := runtime.Caller(3); ok {
			entry.File = file
			entry.Line = line
			if fn := runtime.FuncForPC(pc); fn != nil {
				entry.Function = fn.Name()
			}
		}
	}

	// Write log entry
	l.writeEntry(entry)
}

// writeEntry writes the log entry to output
func (l *StructuredLogger) writeEntry(entry LogEntry) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple logging if JSON marshaling fails
		fmt.Fprintf(l.output, "[%s] %s: %s\n", 
			entry.Timestamp.Format(time.RFC3339), 
			entry.Level, 
			entry.Message)
		return
	}

	fmt.Fprintln(l.output, string(jsonData))
}

// Debug logs a debug message
func (l *StructuredLogger) Debug(msg string, fields ...Field) {
	l.log(DEBUG, msg, fields...)
}

// Info logs an info message
func (l *StructuredLogger) Info(msg string, fields ...Field) {
	l.log(INFO, msg, fields...)
}

// Warn logs a warning message
func (l *StructuredLogger) Warn(msg string, fields ...Field) {
	l.log(WARN, msg, fields...)
}

// Error logs an error message
func (l *StructuredLogger) Error(msg string, fields ...Field) {
	l.log(ERROR, msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *StructuredLogger) Fatal(msg string, fields ...Field) {
	l.log(FATAL, msg, fields...)
	os.Exit(1)
}

// WithFields creates a new logger with additional fields
func (l *StructuredLogger) WithFields(fields ...Field) Logger {
	newLogger := &StructuredLogger{
		level:   l.level,
		output:  l.output,
		fields:  make(map[string]interface{}),
		context: l.context,
		caller:  l.caller,
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for _, field := range fields {
		newLogger.fields[field.Key] = field.Value
	}

	return newLogger
}

// WithError creates a new logger with error field
func (l *StructuredLogger) WithError(err error) Logger {
	return l.WithFields(NewField("error", err.Error()))
}

// WithContext creates a new logger with context
func (l *StructuredLogger) WithContext(ctx string) Logger {
	newLogger := &StructuredLogger{
		level:   l.level,
		output:  l.output,
		fields:  make(map[string]interface{}),
		context: ctx,
		caller:  l.caller,
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// SimpleLogger provides a simple logging interface for backward compatibility
type SimpleLogger struct {
	logger *log.Logger
	level  LogLevel
}

// NewSimpleLogger creates a simple logger
func NewSimpleLogger(level LogLevel) *SimpleLogger {
	return &SimpleLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  level,
	}
}

// Debug logs a debug message
func (l *SimpleLogger) Debug(msg string, fields ...Field) {
	if DEBUG >= l.level {
		l.logger.Printf("[DEBUG] %s %s", msg, l.formatFields(fields...))
	}
}

// Info logs an info message
func (l *SimpleLogger) Info(msg string, fields ...Field) {
	if INFO >= l.level {
		l.logger.Printf("[INFO] %s %s", msg, l.formatFields(fields...))
	}
}

// Warn logs a warning message
func (l *SimpleLogger) Warn(msg string, fields ...Field) {
	if WARN >= l.level {
		l.logger.Printf("[WARN] %s %s", msg, l.formatFields(fields...))
	}
}

// Error logs an error message
func (l *SimpleLogger) Error(msg string, fields ...Field) {
	if ERROR >= l.level {
		l.logger.Printf("[ERROR] %s %s", msg, l.formatFields(fields...))
	}
}

// Fatal logs a fatal message and exits
func (l *SimpleLogger) Fatal(msg string, fields ...Field) {
	if FATAL >= l.level {
		l.logger.Printf("[FATAL] %s %s", msg, l.formatFields(fields...))
		os.Exit(1)
	}
}

// WithFields creates a new logger with additional fields
func (l *SimpleLogger) WithFields(fields ...Field) Logger {
	return &SimpleLogger{
		logger: l.logger,
		level:  l.level,
	}
}

// WithError creates a new logger with error field
func (l *SimpleLogger) WithError(err error) Logger {
	return l.WithFields(NewField("error", err.Error()))
}

// WithContext creates a new logger with context
func (l *SimpleLogger) WithContext(ctx string) Logger {
	return &SimpleLogger{
		logger: l.logger,
		level:  l.level,
	}
}

// formatFields formats fields for simple logging
func (l *SimpleLogger) formatFields(fields ...Field) string {
	if len(fields) == 0 {
		return ""
	}

	var parts []string
	for _, field := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", field.Key, field.Value))
	}

	return "[" + strings.Join(parts, " ") + "]"
}

// Global logger instance
var globalLogger Logger = NewDefaultLogger()

// SetGlobalLogger sets the global logger
func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}

// GetGlobalLogger returns the global logger
func GetGlobalLogger() Logger {
	return globalLogger
}

// Convenience functions for global logger
func Debug(msg string, fields ...Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	globalLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	globalLogger.Fatal(msg, fields...)
}
