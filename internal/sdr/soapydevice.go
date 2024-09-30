package sdr

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// SoapyDevice provides the go-soapy-sdr interface via dependency injection.
type SoapyDevice struct{}

// Enumerate returns a list of available devices on the system.
func (sD SoapyDevice) Enumerate(args map[string]string) []map[string]string {
	return device.Enumerate(args)
}
