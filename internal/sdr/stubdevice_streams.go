package sdr

import (
	"errors"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

var stream *StreamCS8

func (dev *StubDevice) SetupCS8Stream(direction device.Direction,
	channels []uint,
	args map[string]string) (*StreamCS8, error) {
	// Device serial number is used to determine whether to return a valid
	// StreamCS8 or an error.
	switch dev.Args["serial"] {
	case "1":
		stream = nil
		return stream, errors.New("Bad args passed to SetupCS8Stream")
	default:
		// For test purposes, we are only interested that the stream exists, not
		// it's specific values.
		stream = &StreamCS8{&device.SDRStreamCS8{}}
		return stream, nil
	}
}
