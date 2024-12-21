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
	assert.Equal(t, "could not enable Agc", err.Error())
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
	assert.Equal(t, 2, len(elements))

	assert.Contains(t, elements, "RF")
	assert.Contains(t, elements, "IF")
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
	assert.Equal(t, "requested overall gain = 50.1 dB, but must be between 0.0 and 50.0 dB", err.Error())
}

func TestSetOverallGain_Negative(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	// attempting to set overall gain to a negative value.
	err := sdr.SetOverallGain(&stub, testLogger, -2.0)
	assert.NotNil(t, err)
	assert.Equal(t, "requested overall gain = -2.0 dB, but must be between 0.0 and 50.0 dB", err.Error())

}

func TestSetOverallGain(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	err := sdr.SetOverallGain(&stub, testLogger, 40.)
	assert.Nil(t, err)
	gain := sdr.GetOverallGain(&stub, testLogger)
	assert.Equal(t, 40.0, gain)
}

func TestGetElementGain(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	gain, err := sdr.GetElementGain(&stub, testLogger, "RF")
	assert.Nil(t, err)
	assert.Equal(t, 25., gain)
}

func TestGetElementGain_InvalidElement(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	gain, err := sdr.GetElementGain(&stub, testLogger, "Audio")
	assert.NotNil(t, err)
	assert.Equal(t, "gain element 'Audio' is invalid", err.Error())
	assert.Equal(t, 0.0, gain)
}

func TestGetGainElementRange(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	rfRange, err := sdr.GetElementGainRange(&stub, testLogger, "RF")
	assert.Nil(t, err)
	assert.Equal(t, 0.0, rfRange.Minimum)
	assert.Equal(t, 25.0, rfRange.Maximum)
	assert.Equal(t, 1.0, rfRange.Step)
}

func TestGetGainElementRange_BadElement(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	_, err := sdr.GetElementGainRange(&stub, testLogger, "Audio")
	assert.NotNil(t, err)
	assert.Equal(t, "Gain element name: Audio is invalid\n", err.Error())
}

func TestSetElementGain(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	err := sdr.SetElementGain(&stub, testLogger, "RF", 22.0)
	assert.Nil(t, err)
	gain, _ := sdr.GetElementGain(&stub, testLogger, "RF")
	assert.Equal(t, 22.0, gain)

}

func TestSetElementGain_InvalidElement(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	err := sdr.SetElementGain(&stub, testLogger, "Audio", 22.0)
	assert.NotNil(t, err)
	assert.Equal(t, "cannot set gain for non-existent gain element: Audio", err.Error())
}

func TestSetElementGain_InvalidValue(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	err := sdr.SetElementGain(&stub, testLogger, "RF", -1.0)
	assert.NotNil(t, err)
	assert.Equal(t, "cannot set gain for element: RF to -1.0. Requested gain is outside the allowable range: 0.0 to 25.0",
		err.Error())

	err = sdr.SetElementGain(&stub, testLogger, "RF", 25.1)
	assert.NotNil(t, err)
	assert.Equal(t, "cannot set gain for element: RF to 25.1. Requested gain is outside the allowable range: 0.0 to 25.0",
		err.Error())
}
