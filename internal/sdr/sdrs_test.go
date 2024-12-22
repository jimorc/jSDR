package sdr_test

import (
	"testing"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"
	"github.com/stretchr/testify/assert"
)

func TestEnumerateSdrsWithoutAudio(t *testing.T) {
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

	allSdrs := sdr.EnumerateSdrsWithoutAudio(&stub, testLogger)
	assert.Equal(t, 2, len(allSdrs.DevicesMap))
	for k := range allSdrs.DevicesMap {
		assert.NotEqual(t, "Built-in Audio", k)
	}
}

func TestEnumerateSdrsWithoutAudio_NoOtherDevices(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Devices: []map[string]string{
		{"default_output": "False",
			"device_id":     "0",
			"driver":        "audio",
			"label":         "Built-in Audio",
			"default_input": "False"},
	},
	}
	sdrs := sdr.EnumerateSdrsWithoutAudio(&stub, testLogger)
	assert.Equal(t, 0, len(sdrs.DevicesMap))
}

func TestEnumerateSdrsWithoutAudio_NoAudio(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	sdrs := sdr.EnumerateSdrsWithoutAudio(&stub, testLogger)
	assert.Equal(t, 0, len(sdrs.DevicesMap))
}
