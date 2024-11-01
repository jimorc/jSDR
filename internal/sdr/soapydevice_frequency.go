package sdr

import "github.com/pothosware/go-soapy-sdr/pkg/device"

// GetFrequencyRanges retrieves the slice of frequency ranges for the specified direction and channel
// of the device.
func (sD *SoapyDevice) GetFrequencyRanges(direction device.Direction, channel uint) []device.SDRRange {
	return sD.Device.Device.GetFrequencyRange(direction, channel)
}

// GetTunableElements retrieves the list of tunable elements bu name for the device.
//
// The list of tunable elements is expected to be in order from RF to baseband.
func (sD *SoapyDevice) GetTunableElements(direction device.Direction, channel uint) []string {
	return sD.Device.Device.ListFrequencies(direction, channel)
}

// GetTunableElementFrequencyRanges retrieves a slice of frequency ranges for the specified tunable element.
func (sD *SoapyDevice) GetTunableElementFrequencyRanges(direction device.Direction, channel uint, name string) []device.SDRRange {
	return sD.Device.Device.GetFrequencyRangeComponent(direction, channel, name)
}
