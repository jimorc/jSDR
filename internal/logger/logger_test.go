package logger_test

import (
	"strings"
	"testing"

	"github.com/jimorc/jsdr/internal/logger"
)

func TestNew(t *testing.T) {
	logBuf := new(strings.Builder)
	_ = logger.New(logBuf)
}
