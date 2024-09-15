package logger_test

import (
	"strings"
	"testing"
	"time"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestLog_UnformattedMessages(t *testing.T) {
	logBuf := new(strings.Builder)
	l := logger.New(logBuf)

	m := logger.NewLogMessage(logger.Error, "An error message\n")
	l.Log(m)
	// wait for the logging to complete
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, "[Error]: An error message\n", logBuf.String())

	l.Log(logger.NewLogMessage(logger.Info, "An Info message"))
	// wait for the logging to complete
	time.Sleep(10 * time.Millisecond)
	l.Close()

	assert.Equal(t, "[Error]: An error message\n[Info]: An Info message", logBuf.String())
}

func TestLog_FormattedMessages(t *testing.T) {
	logBuf := new(strings.Builder)
	l := logger.New(logBuf)

	l.Log(logger.NewLogMessageWithFormat(logger.Fatal, "Test message with variable: %d\n", 16))
	// wait for the logging to complete
	time.Sleep(10 * time.Millisecond)
	l.Close()

	assert.Equal(t, "[Fatal]: Test message with variable: 16\n", logBuf.String())

	logBuf = new(strings.Builder)
	l = logger.New(logBuf)

	l.Log(logger.NewLogMessageWithFormat(logger.Info, "Test msg with two variables: %d, %s", 4, "str"))
	// wait for the logging to complete
	time.Sleep(10 * time.Millisecond)
	l.Close()

	assert.Equal(t, "[Info]: Test msg with two variables: 4, str", logBuf.String())
}

func TestDefaultLoggingLevel(t *testing.T) {
	logBuf := new(strings.Builder)
	l := logger.New(logBuf)

	// Default level is Info, so these messages should be logged.
	l.Log(logger.NewLogMessage(logger.Info, "Info message 1"))
	// wait for the logging to complete
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, "[Info]: Info message 1", logBuf.String())

	// Default logging level is Info, so these messages should not be logged.
	l.Log(logger.NewLogMessage(logger.Debug, "Debug message 1"))
	l.Log(logger.NewLogMessageWithFormat(logger.Debug, "Debug message %d", 2))

	// wait for the logging to complete
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, "[Info]: Info message 1", logBuf.String())

	// But this message should.
	l.Log(logger.NewLogMessage(logger.Error, "An error message"))
	// wait for the logging to complete
	time.Sleep(10 * time.Millisecond)
	l.Close()
	assert.Equal(t, "[Info]: Info message 1[Error]: An error message", logBuf.String())
}

func TestSetLoggingLevel(t *testing.T) {
	logBuf := new(strings.Builder)
	l := logger.New(logBuf)

	// Default level is Info, so these messages should be logged.
	l.Log(logger.NewLogMessage(logger.Info, "Info message 1"))
	// wait for the logging to complete
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, "[Info]: Info message 1", logBuf.String())

	l.SetMaxLevel(logger.Error)

	// Default logging level is Info, so these messages should not be logged.
	l.Log(logger.NewLogMessage(logger.Info, "Info message 4"))
	l.Log(logger.NewLogMessageWithFormat(logger.Error, "Error message %d", 1))
	l.Log(logger.NewLogMessage(logger.Fatal, "Fatal message"))
	// wait for the logging to complete
	time.Sleep(10 * time.Millisecond)
	l.Close()

	assert.Equal(t, "[Info]: Info message 1[Error]: Error message 1[Fatal]: Fatal message", logBuf.String())
}
