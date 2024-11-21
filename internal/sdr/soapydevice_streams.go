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

// ReadCS8Stream reads numElemsToRead elements from the stream.
//
// Params:
//   - buff: the buffer that stream data from a single channel is stored in. This buffer must be
//
// initialized to [1][2 * MTU] in size before Read is called.
//   - numElemsToRead: The number of elements to read. Because the stream is complex, the actual
//
// number of integers read is 2 * MTU.
//   - outputFlags: The flag indicators of the result of the read operation.
//   - timeoutUs: the timeout in microseconds.
//
// Returns:
//   - timeNs: the timestamp for the data in buff.
//   - numElemsRead: the number of elements read. Since this stream is complex, the total number
//
// of integers returned in buff is 2 times this number.
//   - err: the error if the read is not successful, or nil if the read is successful. On error, the
//
// contents of buff, timeNs, and numElemsRead may not be valid.
func (sD *SoapyDevice) ReadCS8Stream(stream *StreamCS8, buff [][]int, numElemsToRead uint, outputFlags *[1]int, timeoutNs uint) (
	timeNs uint, numElemsRead uint, err error) {
	return 0, 0, nil
}
