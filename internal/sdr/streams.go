package sdr

import (
	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// CS8Streams defines the interface for StreamCS8 streams
type CS8Streams interface {
	SetupCS8Stream(device.Direction, []uint, map[string]string) (*StreamCS8, error)
	CloseCS8Stream(*StreamCS8) error
}

/*type CS8StreamFuncs interface {
	Close() error
}*/

// StreamCS8 is the stream for CS8 data.
type StreamCS8 struct {
	stream *device.SDRStreamCS8
	device CS8Streams
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
