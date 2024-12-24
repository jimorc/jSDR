package sdrdevice_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdrdevice"

	"github.com/stretchr/testify/assert"
)

func TestNewFromPreferences(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	app.NewWithID("com.github.jimorc.jsdrtestnew")
	sdrPrefs := sdrdevice.NewFromPreferences(testLogger)
	sdrPrefs.Device = "01"
	sdrPrefs.SampleRate = "256 kHz"
	sdrPrefs.Antenna = "D"
	sdrPrefs.SavePreferences(testLogger)
	sdrPrefs2 := sdrdevice.NewFromPreferences(testLogger)
	assert.Equal(t, "01", sdrPrefs2.Device)
	assert.Equal(t, "256 kHz", sdrPrefs2.SampleRate)
	assert.Equal(t, "D", sdrPrefs2.Antenna)
}

func TestClearPreferences(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	app.NewWithID("com.github.jimorc.jsdrtest")
	sdrPrefs := sdrdevice.NewFromPreferences(testLogger)
	sdrPrefs.Device = "01"
	sdrPrefs.SampleRate = "256 kHz"
	sdrPrefs.Antenna = "D"

	sdrPrefs.SavePreferences(testLogger)
	sdrPrefs.ClearPreferences(testLogger)
	assert.Equal(t, "", sdrPrefs.Device)
	assert.Equal(t, "", sdrPrefs.SampleRate)
	assert.Equal(t, "", sdrPrefs.Antenna)
	assert.Equal(t, "", fyne.CurrentApp().Preferences().String("device"))
	assert.Equal(t, "", fyne.CurrentApp().Preferences().String("samplerate"))
	assert.Equal(t, "", fyne.CurrentApp().Preferences().String("antenna"))
}
