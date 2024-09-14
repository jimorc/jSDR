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
func (l *Logger) Log(m *LogMessage) {
	if l.level < m.level {
		return
	}
	logMsg := ""
	if m.format == "" {
		logMsg = m.message
	} else {
		logMsg = fmt.Sprintf(m.format, m.args...)
	}
	msg := fmt.Sprintf("[%s]: %s", levelAsString(m.level), logMsg)
	l.writer.WriteString(msg)
}

// SetMaxLevel sets the max logging level.
func (l *Logger) SetMaxLevel(level LoggingLevel) {
	l.level = level
}

func levelAsString(level LoggingLevel) string {
	if level < Fatal || level > Debug {
		return fmt.Sprintf("%s:%d", LevelsAsStrings[0], level)
	}
	return LevelsAsStrings[level]
}

// LogMessage contains the information needed to generate a log message.
type LogMessage struct {
	level   LoggingLevel
	message string
	format  string
	args    []any
}

// NewLogMessage creates a log message without formatting. It is of the form:
// "[level]: msg"
func NewLogMessage(level LoggingLevel, msg string) *LogMessage {
	return &LogMessage{level: level, message: msg}
}

// NewLogMessageWithFormat creates a log message with formatting info and arguments.
func NewLogMessageWithFormat(level LoggingLevel, fmt string, args ...any) *LogMessage {
	return &LogMessage{level: level, format: fmt, args: args}
}
