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
