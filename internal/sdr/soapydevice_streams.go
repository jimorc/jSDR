package sdr

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// SetupCS8Stream initializes a stream for the specified direction, channels, and args.
func (sD *SoapyDevice) SetupCS8Stream(direction device.Direction,
	channels []uint,
	args map[string]string) (*StreamCS8, error) {
	stream, err := sD.Device.Device.SetupSDRStreamCS8(direction, channels, args)
	return &StreamCS8{stream: stream, device: sD}, err
}

// Close closes the specified stream.
func (sD *SoapyDevice) CloseCS8Stream(stream *StreamCS8) error {
	err := stream.stream.Close()
	stream.stream = nil
	return err
}

// GetCS8MTU returns the CS8 stream's maximum transmission unit.
func (sD *SoapyDevice) GetCS8MTU(stream *StreamCS8) int {
	return stream.stream.GetMTU()
}

// Activate activates the CS8 stream.
func (sD *SoapyDevice) Activate(stream *StreamCS8, flag device.StreamFlag,
	timeNs int, numElems int) error {
	return stream.stream.Activate(flag, timeNs, numElems)
}

// Deactivate deactivates the active stream.
func (sD *SoapyDevice) Deactivate(stream *StreamCS8, flag device.StreamFlag,
	timeNS int) error {
	return stream.stream.Deactivate(flag, timeNS)
}
