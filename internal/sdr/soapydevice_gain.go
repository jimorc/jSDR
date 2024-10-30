package sdr

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// AgcIsEnabled returns whether AGC is enabled for receive channel 0 of the device.
// You should call SupportsAGC to determine if the device supports AGC before calling AgcIsEnabled.
//
// Returns true if AGC is enabled.
func (sD *SoapyDevice) AgcIsEnabled(direction device.Direction, channel uint) bool {
	return sD.Device.Device.GetGainMode(direction, channel)
}

// EnableAgc enables AGC mode for the specified direction and channel.
func (sD *SoapyDevice) EnableAgc(direction device.Direction, channel uint, enable bool) error {
	return sD.Device.Device.SetGainMode(direction, channel, enable)
}

// GetGainElementNames returns a list of names for the gain elements for the specified direction and channel.
func (sD *SoapyDevice) GetGainElementNames(direction device.Direction, channel uint) []string {
	return sD.Device.Device.ListGains(direction, channel)
}

// GetOverallGain returns the overall gain for the specified direction and channel.
func (sD *SoapyDevice) GetOverallGain(direction device.Direction, channel uint) float64 {
	return sD.Device.Device.GetGain(direction, channel)
}

// SetOverallGain sets the overall gain for the specified direction and channel.
//
// The overall gain is distributed automatically across the available elements.
func (sD *SoapyDevice) SetOverallGain(direction device.Direction, channel uint, gain float64) error {
	return sD.Device.Device.SetGain(direction, channel, gain)
}

// GetElementGain gets the gain value in dB for the specified element.
//
// The returned error is always nil.
func (sD *SoapyDevice) GetElementGain(direction device.Direction, channel uint, eltName string) (float64, error) {
	return sD.Device.Device.GetGainElement(direction, channel, eltName), nil
}

// SetElementGain sets the gain value in dB for the specified element.
func (sD *SoapyDevice) SetElementGain(direction device.Direction, channel uint, eltName string, gain float64) error {
	return sD.Device.Device.SetGainElement(direction, channel, eltName, gain)
}

// GetElementGainRange returns the SDRRange for the specified gain element.
func (sD *SoapyDevice) GetElementGainRange(direction device.Direction, channel uint, eltName string) device.SDRRange {
	return sD.Device.Device.GetGainElementRange(direction, channel, eltName)
}
