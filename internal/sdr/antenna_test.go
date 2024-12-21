package sdr_test

import (
	"testing"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"

	"github.com/stretchr/testify/assert"
)

func TestGetAntennaNames(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	antennas := sdr.GetAntennaNames(&stub, testLogger)
	assert.Equal(t, []string{"RX"}, antennas)
}

func TestGetCurrentAntenna(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	ant := sdr.GetCurrentAntenna(&stub, testLogger)
	assert.Equal(t, "RX", ant)
}

func TestSetAntenna(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	err := sdr.SetAntenna(&stub, testLogger, "RX")
	assert.Nil(t, err)
}

func TestSetAntenna_BadName(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	err := sdr.SetAntenna(&stub, testLogger, "RX2")
	assert.NotNil(t, err)
	assert.Equal(t, "invalid antenna: RX2", err.Error())
}
