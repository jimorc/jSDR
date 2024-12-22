package sdr_test

import (
	"testing"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
