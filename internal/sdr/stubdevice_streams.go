package sdr

import (
	"errors"
	"time"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

var cs8EltsRead uint

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
// As StubDevice is a test device, 10000 is returned. This matches the MTU value for
// use in ReadCS8Stream, below
func (dev *StubDevice) GetCS8MTU(stream *StreamCS8) int {
	return 10000
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
	switch dev.Args["serial"] {
	case "4":
		switch cs8EltsRead {
		case 0:
			for i := 0; i < 5000/2; i++ {
				buff[0][4*i] = -2
				buff[0][4*i+1] = 0
				buff[0][4*i+2] = -1
				buff[0][4*i+3] = -2
			}
			cs8EltsRead = 5000
			return uint(time.Now().UTC().Nanosecond()), 5000, nil
		case 5000:
			for i := 5000 / 2; i < 8000/2; i++ {
				buff[0][4*i] = -2
				buff[0][4*i+1] = 0
				buff[0][4*i+2] = -1
				buff[0][4*i+3] = -2
			}
			cs8EltsRead = 8000
			return uint(time.Now().UTC().Nanosecond()), 3000, nil
		case 8000:
			for i := 8000 / 2; i < 10000/2; i++ {
				buff[0][4*i] = -2
				buff[0][4*i+1] = 0
				buff[0][4*i+2] = -1
				buff[0][4*i+3] = -2
			}
			cs8EltsRead = 0
			return uint(time.Now().UTC().Nanosecond()), 2000, nil
		}
	default:
		for i := 0; i < int(numElemsToRead/2); i++ {
			buff[0][4*i] = -2
			buff[0][4*i+1] = 0
			buff[0][4*i+2] = -1
			buff[0][4*i+3] = -2
		}
	}
	return uint(time.Now().UTC().Nanosecond()), numElemsToRead, nil
}
