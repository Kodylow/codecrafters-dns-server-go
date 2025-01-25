package gotracer

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// ANSI escape codes for colors
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// LogLevel is our severity threshold.
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger represents a custom logger with color support and extra features.
type Logger struct {
	Info    *log.Logger
	Warn    *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
	level   LogLevel
	outputs []io.Writer
	mu      sync.Mutex
}

// New creates a new Logger instance with colored output and default outputs.
func New() *Logger {
	flags := log.Ldate | log.Ltime | log.Lshortfile
	return &Logger{
		Info:    log.New(os.Stdout, colorBlue+"INFO: "+colorReset, flags),
		Warn:    log.New(os.Stdout, colorYellow+"WARN: "+colorReset, flags),
		Error:   log.New(os.Stderr, colorRed+"ERROR: "+colorReset, flags),
		Debug:   log.New(os.Stdout, colorGreen+"DEBUG: "+colorReset, flags),
		level:   LevelDebug,
		outputs: []io.Writer{os.Stdout},
	}
}

// SetOutput replaces the default outputs entirely with a single writer.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.outputs = []io.Writer{w}
	l.Info.SetOutput(w)
	l.Warn.SetOutput(w)
	l.Error.SetOutput(w)
	l.Debug.SetOutput(w)
}

// AddOutput appends an additional output writer (e.g., file or socket).
func (l *Logger) AddOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.outputs = append(l.outputs, w)
	multiOut := io.MultiWriter(l.outputs...)
	l.Info.SetOutput(multiOut)
	l.Warn.SetOutput(multiOut)
	l.Error.SetOutput(multiOut)
	l.Debug.SetOutput(multiOut)
}

// SetLevel updates the global logging threshold.
func (l *Logger) SetLevel(newLevel LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = newLevel
}

// Log is a unified entry point for all log levels.
func (l *Logger) Log(ctx context.Context, level LogLevel, msg string, fields map[string]interface{}) {
	l.mu.Lock() // minimal locking, might consider finer-grained or none at all
	defer l.mu.Unlock()

	merged := mergeFields(FieldsFromContext(ctx), fields)
	formattedMsg := buildMessage(msg, merged)

	switch {
	case level >= LevelError && l.level <= LevelError:
		l.Error.Output(2, formattedMsg)
	case level == LevelWarn && l.level <= LevelWarn:
		l.Warn.Output(2, formattedMsg)
	case level == LevelInfo && l.level <= LevelInfo:
		l.Info.Output(2, formattedMsg)
	case level == LevelDebug && l.level <= LevelDebug:
		l.Debug.Output(2, formattedMsg)
	}
}

// Infof logs at the INFO level with optional extra fields.
func (l *Logger) Infof(msg string, fields map[string]interface{}) {
	if l.level <= LevelInfo {
		l.Info.Output(2, buildMessage(msg, fields))
	}
}

// Warnf logs at the WARN level with optional extra fields.
func (l *Logger) Warnf(msg string, fields map[string]interface{}) {
	if l.level <= LevelWarn {
		l.Warn.Output(2, buildMessage(msg, fields))
	}
}

// Errorf logs at the ERROR level with optional extra fields.
func (l *Logger) Errorf(msg string, fields map[string]interface{}) {
	if l.level <= LevelError {
		l.Error.Output(2, buildMessage(msg, fields))
	}
}

// Debugf logs at the DEBUG level with optional extra fields.
func (l *Logger) Debugf(msg string, fields map[string]interface{}) {
	if l.level <= LevelDebug {
		l.Debug.Output(2, buildMessage(msg, fields))
	}
}

// -- NEW CONTEXT-BASED METHODS --

func (l *Logger) InfofCtx(ctx context.Context, msg string, fields map[string]interface{}) {
	if l.level <= LevelInfo {
		merged := mergeFields(FieldsFromContext(ctx), fields)
		l.Info.Output(2, buildMessage(msg, merged))
	}
}

func (l *Logger) WarnfCtx(ctx context.Context, msg string, fields map[string]interface{}) {
	if l.level <= LevelWarn {
		merged := mergeFields(FieldsFromContext(ctx), fields)
		l.Warn.Output(2, buildMessage(msg, merged))
	}
}

func (l *Logger) ErrorfCtx(ctx context.Context, msg string, fields map[string]interface{}) {
	if l.level <= LevelError {
		merged := mergeFields(FieldsFromContext(ctx), fields)
		l.Error.Output(2, buildMessage(msg, merged))
	}
}

func (l *Logger) DebugfCtx(ctx context.Context, msg string, fields map[string]interface{}) {
	if l.level <= LevelDebug {
		merged := mergeFields(FieldsFromContext(ctx), fields)
		l.Debug.Output(2, buildMessage(msg, merged))
	}
}

// buildMessage constructs a string from the message plus any structured fields.
func buildMessage(msg string, fields map[string]interface{}) string {
	var extra string
	for k, v := range fields {
		extra += fmt.Sprintf(" %s=%v", k, v)
	}
	return fmt.Sprintf("%s%s", msg, extra)
}

// mergeFields combines context-fields and call-site fields
func mergeFields(ctxFields, callFields map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for k, v := range ctxFields {
		merged[k] = v
	}
	for k, v := range callFields {
		merged[k] = v
	}
	return merged
}

// SetFlags sets the output flags for all loggers.
func (l *Logger) SetFlags(flags int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Info.SetFlags(flags)
	l.Warn.SetFlags(flags)
	l.Error.SetFlags(flags)
	l.Debug.SetFlags(flags)
}
