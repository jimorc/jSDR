package sdr_test

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testLogger logger.Logger

func TestEnumerateWithoutAudio(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Devices: []map[string]string{
		{
			"default_output": "False",
			"device_id":      "0",
			"driver":         "audio",
			"label":          "Built-in Audio",
			"default_input":  "False",
		},
		{
			"driver":       "rtlsdr",
			"label":        "Generic RTL2832U OEM :: 00000101",
			"manufacturer": "Realtek",
			"product":      "RTL2838UHIDIR",
			"serial":       "00000101",
			"tuner":        "Rafael Micro R820T",
		},
		{
			"driver":       "rtlsdr",
			"label":        "Generic RTL2832U OEM :: 00000102",
			"manufacturer": "Realtek",
			"product":      "RTL2838UHIDIR",
			"serial":       "00000102",
			"tuner":        "Rafael Micro R820T"},
	},
	}

	sdrs := sdr.EnumerateWithoutAudio(&stub, testLogger)
	assert.Equal(t, 2, len(sdrs))
	for k := range sdrs {
		assert.NotEqual(t, "Built-in Audio", k)
	}
}

func TestEnumerateWithoutAudio_NoOtherDevices(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Devices: []map[string]string{
		{"default_output": "False",
			"device_id":     "0",
			"driver":        "audio",
			"label":         "Built-in Audio",
			"default_input": "False"},
	},
	}
	sdrs := sdr.EnumerateWithoutAudio(&stub, testLogger)
	assert.Equal(t, 0, len(sdrs))
}

func TestEnumerateWithoutAudio_NoAudio(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	sdrs := sdr.EnumerateWithoutAudio(&stub, testLogger)
	assert.Equal(t, 0, len(sdrs))
}

func TestMake(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	err := sdr.Make(&stub, map[string]string{
		"driver":       "rtlsdr",
		"label":        "Generic RTL2832U OEM :: 00000102",
		"manufacturer": "Realtek",
		"product":      "RTL2838UHIDIR",
		"serial":       "00000102",
		"tuner":        "Rafael Micro R820T"}, testLogger)
	assert.NotNil(t, stub.Device)
	assert.NotNil(t, stub.Device.Device)
	assert.Nil(t, err)
}

func TestBadMake(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	err := sdr.Make(&stub, map[string]string{}, testLogger)
	assert.Nil(t, stub.Device)
	assert.NotNil(t, err)
	// The following error message is returned from StubDevice only. SoapyDevice would return
	// a different error message in case of error.
	assert.Equal(t, "No arguments provided", err.Error())
}

func TestUnmake(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	err := sdr.Make(&stub, map[string]string{
		"driver":       "rtlsdr",
		"label":        "Generic RTL2832U OEM :: 00000102",
		"manufacturer": "Realtek",
		"product":      "RTL2838UHIDIR",
		"serial":       "00000102",
		"tuner":        "Rafael Micro R820T"}, testLogger)
	require.NotNil(t, stub.Device)
	require.NotNil(t, stub.Device.Device)
	require.Nil(t, err)
	err = sdr.Unmake(&stub, testLogger)
	assert.Nil(t, err)
}

func TestBadUnmake(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	err := sdr.Unmake(&stub, testLogger)
	assert.NotNil(t, err)
}

func TestGetHardwareKey(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	_ = sdr.Make(&stub, map[string]string{
		"driver":       "rtlsdr",
		"label":        "Generic RTL2832U OEM :: 00000102",
		"manufacturer": "Realtek",
		"product":      "RTL2838UHIDIR",
		"serial":       "00000102",
		"tuner":        "Rafael Micro R820T"}, testLogger)
	hwKey := stub.GetHardwareKey()
	assert.Equal(t, "hardKey", hwKey)
}

func TestSetSampleRate_SameAsCurrentRate(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	err := sdr.Make(&stub, map[string]string{
		"driver":       "rtlsdr",
		"label":        "Generic RTL2832U OEM :: 00000102",
		"manufacturer": "Realtek",
		"product":      "RTL2838UHIDIR",
		"serial":       "1",
		"tuner":        "Rafael Micro R820T"}, testLogger)
	require.Nil(t, err)
	re, err := regexp.Compile("\\d+\\.\\d+")
	require.Nil(t, err)
	rate := re.FindString(sdr.GetSampleRate(&stub, testLogger))
	sampleRate, err := strconv.ParseFloat(rate, 64)
	sampleRate *= 1e6
	err = sdr.SetSampleRate(&stub, testLogger, sampleRate)
	assert.Nil(t, err)
}

func TestSetSampleRate_Mismatch(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	err := sdr.Make(&stub, map[string]string{
		"driver":       "rtlsdr",
		"label":        "Generic RTL2832U OEM :: 00000102",
		"manufacturer": "Realtek",
		"product":      "RTL2838UHIDIR",
		"serial":       "0",
		"tuner":        "Rafael Micro R820T"}, testLogger)
	require.Nil(t, err)
	err = sdr.SetSampleRate(&stub, testLogger, 1.024*1e6)
	assert.NotNil(t, err)
	assert.Equal(t, "Attempt to set sample rate to 1024000.0 failed. Sample rate is 2048000.0", err.Error())
}

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
