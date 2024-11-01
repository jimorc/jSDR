package sdr

import "github.com/pothosware/go-soapy-sdr/pkg/device"

// GetFrequencyRanges retrieves the slice of frequency ranges for the specified direction and channel
// of the device.
func (sD *SoapyDevice) GetFrequencyRanges(direction device.Direction, channel uint) []device.SDRRange {
	return sD.Device.Device.GetFrequencyRange(direction, channel)
}

// GetTunableElementNames retrieves the list of tunable elements bu name for the device.
//
// The list of tunable elements is expected to be in order from RF to baseband.
func (sD *SoapyDevice) GetTunableElementNames(direction device.Direction, channel uint) []string {
	return sD.Device.Device.ListFrequencies(direction, channel)
}

// GetTunableElementFrequencyRanges retrieves a slice of frequency ranges for the specified tunable element.
func (sD *SoapyDevice) GetTunableElementFrequencyRanges(direction device.Direction, channel uint, name string) []device.SDRRange {
	return sD.Device.Device.GetFrequencyRangeComponent(direction, channel, name)
}

// GetTunableElementFrequency retrieves the current frequency value in Hz for the tunable element.
func (sD *SoapyDevice) GetTunableElementFrequency(direction device.Direction, channel uint, name string) float64 {
	return sD.Device.Device.GetFrequencyComponent(direction, channel, name)
}

// SetTunableElementFrequency sets the tunable element to the specified value.
//
// For RX, this specifies the down-conversion frequency.
// device.SDRDevice.SetFrequencyComponent, which is called by this method, accepts a map of optional component arguments.
// SetTunableElementFrequency passes no arguments to device.SDRDevice.SetFrequencyComponent.
//
// Returns an error, or nil if successful
func (sD *SoapyDevice) SetTunableElementFrequency(direction device.Direction, channel uint, name string, newFreq float64) error {
	return sD.Device.Device.SetFrequencyComponent(direction, channel, name, newFreq, map[string]string{})
}

// GetOverallCenterFrequency retrieves the overall center frequency of the tunable element chain.
func (sD *SoapyDevice) GetOverallCenterFrequency(direction device.Direction, channel uint) float64 {
	return sD.Device.Device.GetFrequency(direction, channel)
}
