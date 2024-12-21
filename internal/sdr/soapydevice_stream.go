package sdr

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// GetStreamFormats retrieves the formats for the device's specified direction and channel.
func (sD *SoapyDevice) GetStreamFormats(direction device.Direction, channel uint) []string {
	return sD.Device.Device.GetStreamFormats(direction, channel)
}

// GetNativeStreamFormat retrieves the native format and full scale value for the specified channel and direction.
func (sD *SoapyDevice) GetNativeStreamFormat(direction device.Direction, channel uint) (string, float64) {
	return sD.Device.Device.GetNativeStreamFormat(direction, channel)
}
