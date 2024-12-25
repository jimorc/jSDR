package sdr_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMake(t *testing.T) {
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
	assert.NotNil(t, stub.Device)
	assert.NotNil(t, stub.Device.Device)
	assert.Nil(t, err)
}

func TestBadMake(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Device: &sdr.Sdr{}}
	stub.Device.DeviceProperties = map[string]string{}
	err := stub.Device.Make(&stub, testLogger)
	assert.Nil(t, stub.Device.Device)
	assert.NotNil(t, err)
	// The following error message is returned from StubDevice only. SoapyDevice would return
	// a different error message in case of error.
	assert.Equal(t, "no arguments provided", err.Error())
}

func TestUnmake(t *testing.T) {
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
	hwKey := stub.GetHardwareKey()
	assert.Equal(t, "hardKey", hwKey)
}

func TestLoadSavePreferences(t *testing.T) {
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
	app.NewWithID("com.github.jimorc.jsdrtestnew")
	stub.Device.LoadPreferences(testLogger)
	stub.Device.DeviceName = "01"
	stub.Device.SampleRate = 256000.
	stub.Device.Antenna = "D"
	stub.Device.SavePreferences(testLogger)
	stub2 := sdr.StubDevice{Device: &sdr.Sdr{}}
	stub2.Device.DeviceProperties = map[string]string{
		"driver":       "rtlsdr",
		"label":        "Generic RTL2832U OEM :: 00000102",
		"manufacturer": "Realtek",
		"product":      "RTL2838UHIDIR",
		"serial":       "00000102",
		"tuner":        "Rafael Micro R820T"}
	err = stub.Device.Make(&stub, testLogger)
	require.Nil(t, err)
	stub2.Device.LoadPreferences(testLogger)
	assert.Equal(t, "01", stub2.Device.DeviceName)
	assert.Equal(t, 256000., stub2.Device.SampleRate)
	assert.Equal(t, "D", stub2.Device.Antenna)
}

func TestClearPreferences(t *testing.T) {
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
	app.NewWithID("com.github.jimorc.jsdrtest")
	stub.Device.LoadPreferences(testLogger)
	stub.Device.DeviceName = "01"
	stub.Device.SampleRate = 256000.
	stub.Device.Antenna = "D"

	stub.Device.SavePreferences(testLogger)
	stub.Device.ClearPreferences(testLogger)
	assert.Equal(t, "", stub.Device.DeviceName)
	assert.Equal(t, 0., stub.Device.SampleRate)
	assert.Equal(t, "", stub.Device.Antenna)
	assert.Equal(t, "", fyne.CurrentApp().Preferences().String("device"))
	assert.Equal(t, 0., fyne.CurrentApp().Preferences().Float("samplerate"))
	assert.Equal(t, "", fyne.CurrentApp().Preferences().String("antenna"))
}
