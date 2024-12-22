package sdrdevice_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdrdevice"

	"github.com/stretchr/testify/assert"
)

func TestClear(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	app.NewWithID("com.github.jimorc.jsdrtest")
	//	var sdrDevice sdrdevice.SdrDevice
	sdrdevice.ClearDeviceSettings(testLogger)
	assert.Equal(t, "", fyne.CurrentApp().Preferences().String("device"))
	assert.Equal(t, "", fyne.CurrentApp().Preferences().String("samplerate"))
	assert.Equal(t, "", fyne.CurrentApp().Preferences().String("antenna"))
}
