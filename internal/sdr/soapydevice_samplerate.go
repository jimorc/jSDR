package sdr

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// GetSampleRateRange returns a slice of sample rate ranges for the specified device.
func (sD *SoapyDevice) GetSampleRateRange(direction device.Direction, channel uint) []device.SDRRange {
	return sD.Device.Device.GetSampleRateRange(direction, channel)
}

// GetSampleRate returns the currently set sample rate for the device.
// If SetSampleRate has not been called, this is probably the device's default value.
func (sD *SoapyDevice) GetSampleRate(direction device.Direction, channel uint) float64 {
	sD.Device.SampleRate = sD.Device.Device.GetSampleRate(direction, channel)
	return sD.Device.SampleRate
}

func (sD *SoapyDevice) SetSampleRate(direction device.Direction, channel uint, rate float64) error {
	return sD.Device.Device.SetSampleRate(direction, channel, rate)
}
