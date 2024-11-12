package sdr

import (
	"errors"
	"time"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

func (dev *StubDevice) SetupCS8Stream(direction device.Direction,
	channels []uint,
	args map[string]string) (*StreamCS8, error) {
	// Device serial number is used to determine whether to return a valid
	// StreamCS8 or an error.
	switch dev.Args["serial"] {
	case "1":
		return nil, errors.New("Bad args passed to SetupCS8Stream")
	default:
		// For test purposes, we are only interested that the stream exists, not
		// it's specific values.
		stream := &StreamCS8{stream: &device.SDRStreamCS8{}, device: dev}
		return stream, nil
	}
}

// CloseCS8Stream closes the specified test stream.
//
// There is no actual stream, so just a bit of cleanup is done.
func (dev *StubDevice) CloseCS8Stream(stream *StreamCS8) error {
	stream.stream = nil
	return nil
}

// GetCS8MTU returns the stream's maximum transmission unit in number of elements.
//
// As StubDevice is a test device, the value for an RTL_SDR device is returned.
func (dev *StubDevice) GetCS8MTU(stream *StreamCS8) int {
	return 131072
}

// Activate the specified stream. Since StubDevice is a test device, there
// is not much activation to be done.
func (dev *StubDevice) Activate(stream *StreamCS8, flag device.StreamFlag,
	timeNs int, numElems int) error {
	return nil
}

// Deactivate deactivates the specified stream. Since StubDevice is a test
// device, there is not much to be done.
func (dev *StubDevice) Deactivate(stream *StreamCS8, flag device.StreamFlag,
	timeNs int) error {
	switch dev.Args["serial"] {
	case "2":
		return errors.New("Bad device")
	default:
		return nil
	}
}

func (dev *StubDevice) ReadCS8Stream(stream *StreamCS8, buff [][]int, numElemsToRead uint, outputFlags [1]int, timeoutUs uint) (
	timeNs uint, numElemsRead uint, err error) {
	for i := 0; i < int(numElemsToRead/2); i = i + 4 {
		buff[0][4*i] = -2
		buff[0][4*i+1] = 0
		buff[0][4*i+2] = -1
		buff[0][4*i+3] = -2
	}
	return uint(time.Now().UTC().Nanosecond()), numElemsToRead, nil
}
