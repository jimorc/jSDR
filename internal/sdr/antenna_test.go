package sdr_test

import (
	"testing"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"

	"github.com/stretchr/testify/assert"
)

func TestListAntennas(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	antennas := sdr.ListAntennas(&stub, testLogger)
	assert.Equal(t, []string{"RX"}, antennas)
}

func TestGetCurrentAntenna(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	ant := sdr.GetCurrentAntenna(&stub, testLogger)
	assert.Equal(t, "RX", ant)
}
