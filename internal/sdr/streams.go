package sdr

import (
	"errors"
	"time"

	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// CS8Streams defines the interface for StreamCS8 streams
type CS8Streams interface {
	SetupCS8Stream(device.Direction, []uint, map[string]string) (*StreamCS8, error)
	CloseCS8Stream(*StreamCS8) error
	GetCS8MTU(*StreamCS8) int
	Activate(*StreamCS8, device.StreamFlag, int, int) error
	Deactivate(*StreamCS8, device.StreamFlag, int) error
	ReadCS8Stream(*StreamCS8, [][]int, uint, *[1]int, uint) (uint, uint, error)
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
func (stream *StreamCS8) GetMTU(log *logger.Logger) uint {
	mtu := stream.device.GetCS8MTU(stream)
	log.Logf(logger.Debug, "CS8 stream MTU is %d\n", mtu)
	return uint(mtu)
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

// Deactivate deactivates a stream.
//
// Call deactivate when not using using read/write(). The implementation control switches or halt data flow.
//
// Params:
//   - flags: optional flag indicators about the stream. Not all implementations will support the full range of options.
//     In this case, the implementation returns ErrorNotSupported.
//   - timeNs: optional deactivation time in nanoseconds. The timeNs is only valid when the flags have StreamFlagHasTime.
//
// Returns an error or nil in case of success
func (stream *StreamCS8) Deactivate(log *logger.Logger, flags device.StreamFlag, timeNs int) error {
	if !stream.active {
		log.Log(logger.Error, "Attempting to deactivate a stream that is not active.\n")
		return errors.New("Attempting to deactivate a stream that is not active")
	}
	err := stream.device.Deactivate(stream, flags, timeNs)
	if err != nil {
		log.Logf(logger.Error, "Error encountered deactivating a stream: %s\n", err.Error())
		return err
	}
	stream.active = false
	return nil
}

// ReadCS8FromStream reads MTU items from the stream. Since the format is CS8, 2 * MTU integers values are read.
//
// Data is read only from RX channel 0.
//
// Params:
//   - buff: an array of buffers that will hold the data that is read. buff must be sized to [1][2*mtu] in the
//
// code that calls ReadCS8FromStream.
//
//   - outputFlags: The flag indicators of the result. Since data is read only from RX channel 0, this array is
//
// one element in size.
//
//   - timeoutUs: the timeout time in microseconds.
//
// Returns:
//   - timeNs: the buffer's timestamp in nanoseconds.
//   - numElemsRead: the number of elements read. This should match the stream's MTU.
//   - err: error, or nil if the call is successful. On error, buff, numElemsRead, and timeNs may not be valid.
func (stream *StreamCS8) ReadCS8FromStream(log *logger.Logger, buff [][]int, elementsToRead uint, outputFlags *[1]int, timeoutUs uint) (
	timeNs uint, numElemsRead uint, err error) {
	if !stream.active {
		log.Log(logger.Error, "Attempting to read from an inactive stream.\n")
		return 0, 0, errors.New("Attempting to read from an inactive stream")
	}
	var elemsRead uint
	mtu := stream.GetMTU(log)
	cs8Buff := make([][]int, 1)
	cs8Buff[0] = make([]int, 2*mtu)
	start := time.Now()
	for {
		if numElemsRead < elementsToRead {
			timeNs, elemsRead, err = stream.device.ReadCS8Stream(stream, cs8Buff, elementsToRead, outputFlags, timeoutUs)
			if err != nil {
				log.Logf(logger.Error, "Error encountered while reading CS8 data: %s\n", err.Error())
				return timeNs, numElemsRead, err
			}
			if outputFlags[0] != 0 {
				log.Logf(logger.Debug, "Flags = %d\n", outputFlags[0])
			}
			log.Logf(logger.Debug, "Elements Read: %d\n", elemsRead)
			// for loop used to transfer data because it is 20x to 25x as fast as append.
			for i := uint(0); i < 2*elemsRead; i++ {
				buff[0][i+2*numElemsRead] = cs8Buff[0][i]
			}
			numElemsRead += elemsRead
		} else {
			end := time.Now()
			diff := end.Sub(start).Microseconds()
			log.Logf(logger.Debug, "Time to read CS8 data: %d μs\n", diff)
			return timeNs, numElemsRead, nil
		}
	}
}

// ReadStreamAsCF64Data reads MTU CS8 items from the stream. Since the stream format is CS8, 2 * MTU integers values are read.
// These are then converted to an array of float64 values.
//
// Data is read only from RX channel 0.
//
// Params:
//   - cf64: an array that will hold the data that is read. cf64 must be sized to [2*mtu] in
//
// the code that calls ReadSttreamAsCF64Data.
//
//   - elemsToRead: This should match the stream's MTU.
//
//   - outputFlags: The flag indicators of the result.
//
//   - timeoutUs: the timeout time in microseconds.
//
// Returns:
//
//   - timeNs: the buffer's timestamp in nanoseconds.
//
//   - numElemsRead: the number of elements read. This should match the stream's MTU.
//
//   - err: error, or nil if the call is successful. On error, cf64, numElemsRead, and timeNs are not valid.
//
// The error has already been logged, so it is not necessary to do so again.
func (stream *StreamCS8) ReadStreamAsCF64Data(log *logger.Logger, cf64 []float64, elementsToRead uint,
	outputFlags *int, timeoutUs uint) (
	timeNs uint, numElemsRead uint, err error) {
	mtu := stream.GetMTU(log)
	cs8 := make([][]int, 1)
	cs8[0] = make([]int, 2*mtu)
	var flags [1]int
	timeNs, numElemsRead, err = stream.ReadCS8FromStream(log, cs8, elementsToRead, &flags, timeoutUs)
	if err != nil {
		return 0, 0, err
	}
	*outputFlags = flags[0]
	start := time.Now()
	size := len(cs8[0])
	//	cf64 = make([]float64, size)
	for i := 0; i < size; i++ {
		cf64[i] = float64(cs8[0][i])
	}
	end := time.Now()
	log.Logf(logger.Debug, "Time to convert CS8 data to CF128 data: %d μs\n", end.Sub(start).Microseconds())
	return timeNs, numElemsRead, err
}
