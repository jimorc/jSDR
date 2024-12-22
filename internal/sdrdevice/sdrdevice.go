package sdrdevice

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"
)

// SdrDevice holds SDR device properties across jsdr app executions.
type SdrDevice struct {
	// Name of the current device
	Device string
}

// LoadFromApp loads the values in the SdrDevice struct from those last saved by the app.
//
// Params:
//
//	log is the logger to write messages to.
func (s *SdrDevice) LoadFromApp(log *logger.Logger) {
	var msg strings.Builder
	msg.WriteString("Loading SDRDevice data from the app:\n")
	s.Device = fyne.CurrentApp().Preferences().String("device")
	msg.WriteString(fmt.Sprintf("         Device: %s\n", s.Device))

	log.Log(logger.Debug, msg.String())

	sdrs := sdr.EnumerateSdrsWithoutAudio(sdr.SoapyDev, log)
	if !sdrs.Contains(s.Device, log) {
		s.Clear(log)
	}
}

// Clear clears all values in the SdrDevice struct. This should only be called if the
// previously stored SDR is no longer connected to the computer.
//
// Params:
//
//	log is the logger to write messages to.
func (s *SdrDevice) Clear(log *logger.Logger) {
	s.Device = ""
	log.Log(logger.Debug, "SdrDevice has been cleared\n")
}
