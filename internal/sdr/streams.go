package sdr

import (
	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// CS8Streams defines the interface for StreamCS8 streams
type CS8Streams interface {
	SetupCS8Stream(device.Direction, []uint, map[string]string) (*StreamCS8, error)
	CloseCS8Stream(*StreamCS8) error
	GetCS8MTU(*StreamCS8) int
	Activate(*StreamCS8, device.StreamFlag, int, int) error
}

// StreamCS8 is the stream for CS8 data.
type StreamCS8 struct {
	stream *device.SDRStreamCS8
	device CS8Streams
	active bool
}

// SetupCS8Stream initializes a stream for RX channel 0.
//
// All stream API calls should be usable with the new stream object
// after SetupSDRStreamCU8() is complete, regardless of the activity state.
//
// Returns a stream pointer and an error. The returned stream may not be used
// concurrently on multiple go routines.
func SetupCS8Stream(sdrD CS8Streams, log *logger.Logger) (*StreamCS8, error) {
	// TODO: Determine what the "WIRE" value should be. The SoapySDR documentation does not
	// give any specific values, just says 'format of the samples between device and host.
	// I am guessing that means "CS8" here.
	stream, err := sdrD.SetupCS8Stream(device.DirectionRX, []uint{0}, map[string]string{"WIRE": "CS8"})
	if err != nil {
		log.Logf(logger.Error, "Could not set up stream: %s\n", err.Error())
		return nil, err
	}
	log.Log(logger.Debug, "CS8 stream setup complete.\n")
	return stream, err
}

// CloseCS8Stream closes an open CS8 stream, that is, a stream that was set up with a call to
// sdr.SetupCS8Stream
func (stream *StreamCS8) Close(log *logger.Logger) error {
	err := stream.device.CloseCS8Stream(stream)
	if err != nil {
		log.Logf(logger.Error, "Could not close a stream: %s\n", err.Error())
		return err
	}
	log.Log(logger.Debug, "Stream closed.\n")
	return nil
}

// GetMTU gets stream's maximum transmission unit in number of elements.
// The MTU specifies the maximum payload transfer in a stream operation. This value can be used as a stream buffer
// allocation size that can best optimize throughput given the underlying stream implementation.
//
// Return the MTU in number of stream elements (never zero)
func (stream *StreamCS8) GetMTU(log *logger.Logger) int {
	mtu := stream.device.GetCS8MTU(stream)
	log.Logf(logger.Debug, "CS8 stream MTU is %d\n", mtu)
	return mtu
}

// Activate activates a stream.
//
// Call activate to prepare a stream before using read/write(). The implementation control switches or stimulate data
// flow.
//
// Params:
//   - flags: optional flag indicators about the stream. The StreamFlagEndBurst flag can signal end on the finite burst.
//     Not all implementations will support the full range of options. In this case, the implementation returns
//     ErrorNotSupported.
//   - timeNs: optional activation time in nanoseconds. The timeNs is only valid when the flags have StreamFlagHasTime.
//   - numElems: optional element count for burst control. The numElems count can be used to request a finite burst size.
//
// Return an error or nil in case of success
func (stream *StreamCS8) Activate(log *logger.Logger, flag device.StreamFlag, timeNs int, numElems int) error {
	err := stream.device.Activate(stream, flag, timeNs, numElems)
	if err != nil {
		log.Logf(logger.Error, "Error attempting to activate CS8 stream: %s\n", err.Error())
		stream.active = false
		return err
	}
	stream.active = true
	return nil
}
