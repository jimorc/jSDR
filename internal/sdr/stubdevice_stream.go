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
