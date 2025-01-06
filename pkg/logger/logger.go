package logger

import (
	"io"
	"log"
)

// ANSI escape codes for colors
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// Logger represents a custom logger with color support
type Logger struct {
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Debug *log.Logger
}

// New creates a new Logger instance with colored output
func New() *Logger {
	flags := log.Ldate | log.Ltime | log.Lshortfile | log.LUTC
	return &Logger{
		Info:  log.New(log.Writer(), colorBlue+"INFO: "+colorReset, flags),
		Warn:  log.New(log.Writer(), colorYellow+"WARN: "+colorReset, flags),
		Error: log.New(log.Writer(), colorRed+"ERROR: "+colorReset, flags),
		Debug: log.New(log.Writer(), colorGreen+"DEBUG: "+colorReset, flags),
	}
}

// SetOutput changes the output destination for all loggers
func (l *Logger) SetOutput(w io.Writer) {
	l.Info.SetOutput(w)
	l.Warn.SetOutput(w)
	l.Error.SetOutput(w)
	l.Debug.SetOutput(w)
}

// SetFlags sets the output flags for all loggers
func (l *Logger) SetFlags(flags int) {
	l.Info.SetFlags(flags)
	l.Warn.SetFlags(flags)
	l.Error.SetFlags(flags)
	l.Debug.SetFlags(flags)
}
