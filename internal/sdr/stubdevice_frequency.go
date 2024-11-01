package sdr

import (
	"errors"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

var tunableFrequencies = map[string]float64{
	"RF": 100000000.0,
}

var overallCenterFrequency = 100000000.

// GetFrequencyRanges retrieves the frequency ranges supported by the device.
func (dev *StubDevice) GetFrequencyRanges(_ device.Direction, _ uint) []device.SDRRange {
	switch dev.Args["serial"] {
	case "2", "4":
		return []device.SDRRange{{Minimum: 0.0, Maximum: 6e+09, Step: 0.0}}
	case "3":
		return []device.SDRRange{{Minimum: 0.0, Maximum: 6e+09, Step: 0.0},
			{Minimum: 6.1e+09, Maximum: 1e+10, Step: 0.0}}
	default:
		return []device.SDRRange{}
	}
}

// GetTunableElementNames returns the list of tunable elements for this device.
func (dev *StubDevice) GetTunableElementNames(_ device.Direction, _ uint) []string {
	return []string{"RF"}
}

// GetTunableElementFrequencyRanges retrieves a slice of frequency ranges for the specified tunable element.
func (dev *StubDevice) GetTunableElementFrequencyRanges(_ device.Direction, _ uint, _ string) []device.SDRRange {
	return []device.SDRRange{{Minimum: 0, Maximum: 6e+09, Step: 0},
		{Minimum: 6.1e+09, Maximum: 1e+10, Step: 0}}
}

// GetTunableElementFrequency retrieves the tuned frequency in Hz for the specified tunable element.
func (dev *StubDevice) GetTunableElementFrequency(_ device.Direction, _ uint, name string) float64 {
	return tunableFrequencies["RF"]
}

// SetTunableElementFrequency sets the frequency for the tunable element to the specified value.
func (dev *StubDevice) SetTunableElementFrequency(_ device.Direction, _ uint, name string, newFreq float64) error {
	tunableFrequencies[name] = newFreq
	return nil
}

// GetOverallCenterFrequency retrieves the overall center frequency for the device.
func (dev *StubDevice) GetOverallCenterFrequency(_ device.Direction, _ uint) float64 {
	return overallCenterFrequency
}

// SetOverallCenterFrequency sets the overall center frequency for the device.
//
// Returns error if serial number for the device is "4". This is purely for testing purposes.
func (dev *StubDevice) SetOverallCenterFrequency(_ device.Direction, _ uint, newFreq float64, _ map[string]string) error {
	switch dev.Args["serial"] {
	case "4":
		return errors.New("Serial # = 4")
	default:
		overallCenterFrequency = newFreq
		return nil
	}
}
