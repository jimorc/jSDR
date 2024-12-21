package sdr_test

import (
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
	assert.Equal(t, "no arguments provided", err.Error())
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
