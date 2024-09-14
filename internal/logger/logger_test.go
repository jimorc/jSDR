package logger_test

import (
	"strings"
	"testing"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	logBuf := new(strings.Builder)
	l := logger.New(logBuf)

	l.Log(logger.Error, "An error message")

	assert.Equal(t, "[Error]: An error message", logBuf.String())
}

func TestLogf(t *testing.T) {
	logBuf := new(strings.Builder)
	l := logger.New(logBuf)

	l.Logf(logger.Fatal, "Test message with variable: %d", 16)

	assert.Equal(t, "[Fatal]: Test message with variable: 16", logBuf.String())

	logBuf = new(strings.Builder)
	l = logger.New(logBuf)

	l.Logf(logger.Info, "Test msg with two variables: %d, %s", 4, "str")

	assert.Equal(t, "[Info]: Test msg with two variables: 4, str", logBuf.String())
}

func TestLogln(t *testing.T) {
	logBuf := new(strings.Builder)
	l := logger.New(logBuf)

	l.Logln(logger.Fatal, "Fatal msg")

	assert.Equal(t, "[Fatal]: Fatal msg\n", logBuf.String())
}

func TestDefaultLoggingLevel(t *testing.T) {
	logBuf := new(strings.Builder)
	l := logger.New(logBuf)

	// Default level is Info, so these messages should be logged.
	l.Log(logger.Info, "Info message 1")
	l.Logf(logger.Info, "Info message %d", 2)
	l.Logln(logger.Info, "Info message 3")

	assert.Equal(t, "[Info]: Info message 1[Info]: Info message 2[Info]: Info message 3\n", logBuf.String())

	// Default logging level is Info, so these messages should not be logged.
	l.Log(logger.Debug, "Debug message 1")
	l.Logf(logger.Debug, "Debug message %d", 2)
	l.Logln(logger.Debug, "Debug message 3")

	assert.Equal(t, "[Info]: Info message 1[Info]: Info message 2[Info]: Info message 3\n", logBuf.String())
}

func TestSetLoggingLevel(t *testing.T) {
	logBuf := new(strings.Builder)
	l := logger.New(logBuf)

	// Default level is Info, so these messages should be logged.
	l.Log(logger.Info, "Info message 1")
	l.Logf(logger.Info, "Info message %d", 2)
	l.Logln(logger.Info, "Info message 3")

	assert.Equal(t, "[Info]: Info message 1[Info]: Info message 2[Info]: Info message 3\n", logBuf.String())

	l.SetMaxLevel(logger.Error)

	// Default logging level is Info, so these messages should not be logged.
	l.Log(logger.Info, "Info message 4")
	l.Logf(logger.Error, "Error message %d", 1)
	l.Logln(logger.Fatal, "Fatal message")

	assert.Equal(t, `[Info]: Info message 1[Info]: Info message 2[Info]: Info message 3
[Error]: Error message 1[Fatal]: Fatal message
`, logBuf.String())
}
