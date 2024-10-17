package sdr_test

import (
	"slices"
	"testing"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSampleRates(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	err := sdr.Make(&stub, map[string]string{
		"driver":       "rtlsdr",
		"label":        "Generic RTL2832U OEM :: 00000102",
		"manufacturer": "Realtek",
		"product":      "RTL2838UHIDIR",
		"serial":       "00000102",
		"tuner":        "Rafael Micro R820T"}, testLogger)
	require.Nil(t, err)
	rates := sdr.GetSampleRates(&stub, testLogger)
	require.Equal(t, 7, len(rates))
	assert.True(t, slices.Index(rates, "0.256 MS/s") != -1)
	assert.True(t, slices.Index(rates, "1.024 MS/s") != -1)
	assert.True(t, slices.Index(rates, "1.6 MS/s") != -1)
	assert.True(t, slices.Index(rates, "2.048 MS/s") != -1)
	assert.True(t, slices.Index(rates, "2.4 MS/s") != -1)
	assert.True(t, slices.Index(rates, "2.8 MS/s") != -1)
	assert.True(t, slices.Index(rates, "3.2 MS/s") != -1)

}

func TestGetSampleRates_NoDevice(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	rates := sdr.GetSampleRates(&stub, testLogger)
	require.Equal(t, 0, len(rates))
}

func TestGetSampleRate(t *testing.T) {
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
	sampleRate := sdr.GetSampleRate(&stub, testLogger)
	assert.Equal(t, "2.048 MS/s", sampleRate)
}

func TestGetSampleRate_BadDevice(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{}
	// stub.Make not called, so no sample rate returned by sdr.GetSampleRate.
	sampleRate := sdr.GetSampleRate(&stub, testLogger)
	assert.Equal(t, "", sampleRate)
}
