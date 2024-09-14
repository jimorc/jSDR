package logger

import (
	"fmt"
	"io"
)

type LoggingLevel uint8

// Logging levels
const (
	_                  = iota
	Fatal LoggingLevel = iota
	Error
	Info
	Debug
)

var LevelsAsStrings [5]string = [5]string{"Undefined", "Fatal", "Error", "Info", "Debug"}

// Logger is a simple logger. It provides a few functions and methods to log information
// to any io.StringWriter.
type Logger struct {
	writer io.StringWriter
	level  LoggingLevel
}

// New creates a new Logger with a max logging level of Info.
func New(writer io.StringWriter) *Logger {
	return &Logger{writer: writer, level: Info}
}

// Log writes a log message prepended by the logging level. A line return is not
// appended to the messsage.
func (l *Logger) Log(level LoggingLevel, messsage string) {
	if l.level < level {
		return
	}
	logMsg := fmt.Sprintf("[%s]: %s", levelAsString(level), messsage)
	l.writer.WriteString(logMsg)
}

// Logf writes a formatted message prepended by the logging level.
func (l *Logger) Logf(level LoggingLevel, format string, args ...any) {
	if l.level < level {
		return
	}
	formatted := fmt.Sprintf(format, args...)
	msg := fmt.Sprintf("[%s]: %s", levelAsString(level), formatted)
	l.writer.WriteString(msg)
}

// Logln writes a log message prepended by the logging level with a
// line return appended.
func (l *Logger) Logln(level LoggingLevel, message string) {
	if l.level < level {
		return
	}
	msg := fmt.Sprintf("%s\n", message)
	l.Log(level, msg)
}

func levelAsString(level LoggingLevel) string {
	if level < Fatal || level > Debug {
		return fmt.Sprintf("%s:%d", LevelsAsStrings[0], level)
	}
	return LevelsAsStrings[level]
}
