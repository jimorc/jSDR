package sdr

import (
	"errors"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// SoapyDevice provides the go-soapy-sdr interface via dependency injection.
type SoapyDevice struct {
	Device *Sdr
}

// Enumerate returns a list of available devices on the system.
func (sD SoapyDevice) Enumerate(args map[string]string) []map[string]string {
	return device.Enumerate(args)
}

// Make makes a new device object for the given device construction args.
// The device pointer will be stored in a table so subsequent calls with the same arguments will produce the same
// device. For every call to make, there should be a matched call to unmake.
//
// Params:
//   - args: device construction key/value argument map
//
// Return a pointer to a new Device object or nil for error
func (sD *SoapyDevice) Make(args map[string]string) error {
	dev, err := device.Make(args)
	if err == nil {
		sD.Device = &Sdr{Device: dev, DeviceProperties: args}
	}
	return err
}

// Unmake unmakes the SDR device if Make previously called for it.
func (sD *SoapyDevice) Unmake() error {
	if sD.Device == nil {
		return errors.New("Attempted to Unmake an SDR that was not successfully created")
	}
	err := sD.Device.Device.Unmake()
	if err == nil {
		sD.Device = nil
	}
	return err
}

// GetHardwareKey returns the hardware key for the specified device.
func (sD SoapyDevice) GetHardwareKey() string {
	return sD.Device.Device.GetHardwareKey()
}

// SupportsAGC returns whether the device supports AGC or not.
//
// Returns true if device supports automatic gain control.
func (sD *SoapyDevice) SupportsAGC(direction device.Direction, channel uint) bool {
	return sD.Device.Device.HasGainMode(direction, channel)
}

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
