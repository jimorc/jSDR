package sdr

import (
	"errors"
	"fmt"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// GetSampleRateRange returns sample rate ranges that match an RTLSDR dongle.
// If the Make has not been called for the device, then an empty slice is returned.
func (dev *StubDevice) GetSampleRateRange(_ device.Direction, _ uint) []device.SDRRange {
	if dev.Device == nil {
		return []device.SDRRange{}
	}
	return []device.SDRRange{{Minimum: 225001, Maximum: 300000},
		{Minimum: 900001, Maximum: 3200000}}
}

// GetSampleRate returns a sample rate.
// If Make has not been called for the StubDevice, then 0.0 is returned. Otherwise, 2000000 is returned.
func (dev *StubDevice) GetSampleRate(_ device.Direction, _ uint) float64 {
	if dev.Device == nil {
		return 0.0
	} else {
		// returned sample rate is based on the device's serial number. This is for testing only!
		switch dev.Args["serial"] {
		case "1":
			rate := 2000000.0
			dev.Device.SampleRate = rate
			return rate
		default:
			return 0.0
		}
	}
}

func (dev *StubDevice) SetSampleRate(_ device.Direction, _ uint, rate float64) error {
	// Err that is returned is based on the device's serial number. This is for testing only!
	// Sets various conditions for testing sdr.GetSampleRate.
	switch dev.Args["serial"] {
	case "0":
		dev.sampleRate = rate
		return errors.New(fmt.Sprintf("Attempt to set sample rate to %.1f failed. Sample rate is 2048000.0", rate))
	case "1":
		dev.sampleRate = rate
		return nil
	}
	return nil
}
