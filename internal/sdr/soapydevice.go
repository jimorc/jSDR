package sdr

import (
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
	sD.Device = &Sdr{Device: dev}
	return err
}

// GetHardwareKey returns the hardware key for the specified device.
func (sD SoapyDevice) GetHardwareKey() string {
	return sD.Device.Device.GetHardwareKey()
}
