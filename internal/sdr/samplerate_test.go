package sdr_test

import (
	"regexp"
	"slices"
	"strconv"
	"testing"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSampleRates(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Device: &sdr.Sdr{}}
	stub.Device.DeviceProperties = map[string]string{
		"driver":       "rtlsdr",
		"label":        "Generic RTL2832U OEM :: 00000102",
		"manufacturer": "Realtek",
		"product":      "RTL2838UHIDIR",
		"serial":       "00000102",
		"tuner":        "Rafael Micro R820T"}
	err := stub.Device.Make(&stub, testLogger)
	require.Nil(t, err)
	rates := stub.Device.GetSampleRates(&stub, testLogger)
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
	stub := sdr.StubDevice{Device: &sdr.Sdr{}}
	stub.Device.DeviceProperties = map[string]string{}
	err := stub.Device.Make(&stub, testLogger)
	require.NotNil(t, err)
	rates := stub.Device.GetSampleRates(&stub, testLogger)
	require.Equal(t, 0, len(rates))
}

func TestGetSampleRate(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Device: &sdr.Sdr{}}
	stub.Device.DeviceProperties = map[string]string{
		"driver":       "rtlsdr",
		"label":        "Generic RTL2832U OEM :: 00000102",
		"manufacturer": "Realtek",
		"product":      "RTL2838UHIDIR",
		"serial":       "1",
		"tuner":        "Rafael Micro R820T"}
	err := stub.Device.Make(&stub, testLogger)
	require.Nil(t, err)
	sampleRate := stub.Device.GetSampleRate(&stub, testLogger)
	assert.Equal(t, "2.048 MS/s", sampleRate)
}

func TestGetSampleRate_BadDevice(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Device: &sdr.Sdr{}}
	// stub.Make not called, so no sample rate returned by sdr.GetSampleRate.
	sampleRate := stub.Device.GetSampleRate(&stub, testLogger)
	assert.Equal(t, "", sampleRate)
}

func TestSetSampleRate_SameAsCurrentRate(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Device: &sdr.Sdr{}}
	stub.Device.DeviceProperties = map[string]string{
		"driver":       "rtlsdr",
		"label":        "Generic RTL2832U OEM :: 00000102",
		"manufacturer": "Realtek",
		"product":      "RTL2838UHIDIR",
		"serial":       "00000102",
		"tuner":        "Rafael Micro R820T"}
	err := stub.Device.Make(&stub, testLogger)
	require.Nil(t, err)
	re, err := regexp.Compile(`\d+\.\d+`)
	require.Nil(t, err)
	rate := re.FindString(stub.Device.GetSampleRate(&stub, testLogger))
	sampleRate, _ := strconv.ParseFloat(rate, 64)
	sampleRate *= 1e6
	err = stub.Device.SetSampleRate(&stub, testLogger, sampleRate)
	assert.Nil(t, err)
}

func TestSetSampleRate_Mismatch(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Device: &sdr.Sdr{}}
	stub.Device.DeviceProperties = map[string]string{
		"driver":       "rtlsdr",
		"label":        "Generic RTL2832U OEM :: 00000102",
		"manufacturer": "Realtek",
		"product":      "RTL2838UHIDIR",
		"serial":       "1",
		"tuner":        "Rafael Micro R820T"}
	err := stub.Device.Make(&stub, testLogger)
	require.Nil(t, err)
	err = stub.Device.SetSampleRate(&stub, testLogger, 1.024*1e6)
	assert.NotNil(t, err)
	assert.Equal(t, "attempt to set sample rate to 1024000.0 failed. Sample rate is 2048000.0", err.Error())
}
