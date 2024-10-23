package sdr_test

import (
	"testing"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"

	"github.com/stretchr/testify/assert"
)

func TestSupportsAGC(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	supportsAGC := sdr.SupportsAGC(&stub, testLogger)
	assert.True(t, supportsAGC)
}

func TestAgcIsEnabled(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	agcEnabled := sdr.AgcIsEnabled(&stub, testLogger)
	assert.True(t, agcEnabled)
}

func TestAgcIsNotEnabled(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	agcEnabled := sdr.AgcIsEnabled(&stub, testLogger)
	assert.False(t, agcEnabled)
}

func TestEnableAgc_Error(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	// serial number of "1" will return an error
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	err := sdr.EnableAgc(&stub, testLogger, true)
	assert.Equal(t, "Could not enable Agc", err.Error())
}

func TestEnableAgc(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	// serial number of "2" will enable Agc
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	err := sdr.EnableAgc(&stub, testLogger, true)
	assert.Nil(t, err)
	assert.True(t, sdr.AgcIsEnabled(&stub, testLogger))

	err = sdr.EnableAgc(&stub, testLogger, false)
	assert.Nil(t, err)
	assert.False(t, sdr.AgcIsEnabled(&stub, testLogger))
}

func TestGetGainElementNames(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	elements := sdr.GetGainElementNames(&stub, testLogger)
	assert.Equal(t, 1, len(elements))
	assert.Equal(t, "RX", elements[0])
}

func TestGetOverallGain(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	gain := sdr.GetOverallGain(&stub, testLogger)
	assert.Equal(t, 50., gain)
}

func TestSetOverallGain_TooLarge(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	// attempting to set gain to 50.1 dB but stub only allows values up to 50.0 dB.
	err := sdr.SetOverallGain(&stub, testLogger, 50.1)
	assert.NotNil(t, err)
	assert.Equal(t, "Requested overall gain = 50.1 dB, but must be between 0. and 50. dB.", err.Error())
}

func TestSetOverallGain_Negative(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	// attempting to set overall gain to a negative value.
	err := sdr.SetOverallGain(&stub, testLogger, -2.0)
	assert.NotNil(t, err)
	assert.Equal(t, "Requested overall gain = -2.0 dB, but must be between 0. and 50. dB.", err.Error())

}

func TestSetOverallGain(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	err := sdr.SetOverallGain(&stub, testLogger, 40.)
	assert.Nil(t, err)
	gain := sdr.GetOverallGain(&stub, testLogger)
	assert.Equal(t, 40.0, gain)
}
