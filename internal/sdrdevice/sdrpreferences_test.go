package sdrdevice_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdrdevice"

	"github.com/stretchr/testify/assert"
)

type TestChanges struct{}

var tC TestChanges

// the following functions allow NewFromPreferences to be called during the tests.
func (TestChanges) SdrChanged()        {}
func (TestChanges) SampleRateChanged() {}
func (TestChanges) AntennaChanged()    {}

func TestNewFromPreferences(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	app.NewWithID("com.github.jimorc.jsdrtestnew")
	sdrPrefs := sdrdevice.NewFromPreferences(tC, testLogger)
	err := sdrPrefs.Device.Set("01")
	assert.Nil(t, err)
	err = sdrPrefs.SampleRate.Set("256 kHz")
	assert.Nil(t, err)
	err = sdrPrefs.Antenna.Set("D")
	assert.Nil(t, err)
	err = sdrPrefs.SavePreferences(testLogger)
	assert.Nil(t, err)
	sdrPrefs2 := sdrdevice.NewFromPreferences(tC, testLogger)
	device, err := sdrPrefs2.Device.Get()
	assert.Nil(t, err)
	assert.Equal(t, "01", device)
	rate, err := sdrPrefs2.SampleRate.Get()
	assert.Nil(t, err)
	assert.Equal(t, "256 kHz", rate)
	antenna, err := sdrPrefs2.Antenna.Get()
	assert.Nil(t, err)
	assert.Equal(t, "D", antenna)
}

func TestClearPreferences(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	app.NewWithID("com.github.jimorc.jsdrtest")
	sdrPrefs := sdrdevice.NewFromPreferences(tC, testLogger)
	err := sdrPrefs.Device.Set("01")
	assert.Nil(t, err)
	err = sdrPrefs.SampleRate.Set("256 kHz")
	assert.Nil(t, err)
	err = sdrPrefs.Antenna.Set("D")
	assert.Nil(t, err)
	err = sdrPrefs.SavePreferences(testLogger)
	assert.Nil(t, err)
	sdrPrefs.ClearPreferences(testLogger)
	value, err := sdrPrefs.Device.Get()
	assert.Nil(t, err)
	assert.Equal(t, "", value)
	value, err = sdrPrefs.SampleRate.Get()
	assert.Nil(t, err)
	assert.Equal(t, "", value)
	value, err = sdrPrefs.Antenna.Get()
	assert.Nil(t, err)
	assert.Equal(t, "", value)
	assert.Equal(t, "", fyne.CurrentApp().Preferences().String("device"))
	assert.Equal(t, "", fyne.CurrentApp().Preferences().String("samplerate"))
	assert.Equal(t, "", fyne.CurrentApp().Preferences().String("antenna"))
}
