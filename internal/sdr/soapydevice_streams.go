package sdr

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// SetupCS8Stream initializes a stream for the specified direction, channels, and args.
func (sD *SoapyDevice) SetupCS8Stream(direction device.Direction,
	channels []uint,
	args map[string]string) (*StreamCS8, error) {
	stream, err := sD.Device.Device.SetupSDRStreamCS8(direction, channels, args)
	return &StreamCS8{stream}, err
}
