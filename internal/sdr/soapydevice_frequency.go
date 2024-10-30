package sdr

import "github.com/pothosware/go-soapy-sdr/pkg/device"

// GetFrequencyRanges retrieves the slice of frequency ranges for the specified direction and channel
// of the device.
func (sD *SoapyDevice) GetFrequencyRanges(direction device.Direction, channel uint) []device.SDRRange {
	return sD.Device.Device.GetFrequencyRange(direction, channel)
}
