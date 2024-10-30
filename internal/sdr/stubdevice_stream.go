package sdr

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// GetStreamFormats returns a slice of stream formats.
//
// The formats returned either match those for RTL_SDR dongles, or no formats.
func (dev *StubDevice) GetStreamFormats(_ device.Direction, _ uint) []string {
	switch dev.Args["serial"] {
	case "2":
		return []string{"CS8", "CS16", "CF32"}
	default:
		return []string{}
	}
}

// GetNativeStreamFormat returns the native stream format and its full scale value.
//
// The returned values match those for an RTL_SDR device.
func (dev *StubDevice) GetNativeStreamFormat(_ device.Direction, _ uint) (string, float64) {
	return "CS8", 0.0
}
