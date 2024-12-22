package sdrdevice

import (
	"fyne.io/fyne/v2"
	"github.com/jimorc/jsdr/internal/logger"
)

// Clear clears all values in the SdrDevice struct. This should only be called if the
// previously stored SDR is no longer connected to the computer.
//
// Params:
//
//	log is the logger to write messages to.
func ClearDeviceSettings(log *logger.Logger) {
	fyne.CurrentApp().Preferences().SetString("device", "")
	fyne.CurrentApp().Preferences().SetString("samplerate", "")
	fyne.CurrentApp().Preferences().SetString("antenna", "")
	log.Log(logger.Debug, "SdrDevice device settings have been cleared\n")
}
