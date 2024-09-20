package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type LoggingLevel uint8

// Logging levels
const (
	None LoggingLevel = iota
	Fatal
	Error
	Info
	Debug
)

var LevelsAsStrings [5]string = [5]string{"Undefined", "Fatal", "Error", "Info", "Debug"}

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

// Logger is a simple logger. It provides a few functions and methods to log information
// to any io.StringWriter.
type Logger struct {
	writer    io.StringWriter
	file      *os.File
	waitGroup sync.WaitGroup
	level     LoggingLevel
	logCh     chan LogMessage
}

// New creates a new Logger with a max logging level of Info.
func New(writer io.StringWriter) *Logger {
	l := &Logger{writer: writer}
	l.waitGroup.Add(1)
	l.logCh = make(chan LogMessage, 100)
	go l.outputMessages()
	l.SetMaxLevel(Info)
	return l
}

// NewFileLogger creates a new Logger that logs to a file.
func NewFileLogger(f string) (*Logger, error) {
	file, err := os.Create(f)
	if err != nil {
		return nil, err
	}
	l := New(file)
	l.file = file
	return l, nil
}

// Close closes the logger.
func (l *Logger) Close() {
	close(l.logCh)
	l.waitGroup.Wait()
	if l.file != nil {
		l.file.Close()
	}
}

// Log queues a log message to be output to the logger. The message is queued only if
// the message level is less than the logger's maximum logging level.
func (l *Logger) Log(m *LogMessage) {
	if l.level < m.level {
		return
	}
	l.logCh <- *m
}

// SetMaxLevel sets the max logging level.
func (l *Logger) SetMaxLevel(level LoggingLevel) {
	l.level = Info
	l.Log(NewLogMessageWithFormat(Info, "Setting max logging level to '%s'\n", levelAsString(level)))
	l.level = level
}

func levelAsString(level LoggingLevel) string {
	if level < Fatal || level > Debug {
		return fmt.Sprintf("%s:%d", LevelsAsStrings[0], level)
	}
	return LevelsAsStrings[level]
}

func (l *Logger) outputMessages() {
	for {
		msg, ok := <-l.logCh
		if !ok {
			l.waitGroup.Done()
			return
		}
		logMsg := ""
		if msg.format == "" {
			logMsg = msg.message
		} else {
			logMsg = fmt.Sprintf(msg.format, msg.args...)
		}
		message := fmt.Sprintf("[%s]: %s", levelAsString(msg.level), logMsg)
		l.writer.WriteString(message)
	}
}
