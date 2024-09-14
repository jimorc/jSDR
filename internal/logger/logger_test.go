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
