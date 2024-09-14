package logger

import (
	"io"
)

// Logger is a simple logger. It provides a few functions and methods to log information
// to any io.Writer.
type Logger struct {
	writer io.Writer
}

// New creates a new Logger.
func New(writer io.Writer) *Logger {
	return &Logger{writer: writer}
}
