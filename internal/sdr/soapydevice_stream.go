package sdr

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// GetStreamFormats retrieves the formats for the device's specified direction and channel.
func (sD *SoapyDevice) GetStreamFormats(direction device.Direction, channel uint) []string {
	return sD.GetStreamFormats(direction, channel)
}
