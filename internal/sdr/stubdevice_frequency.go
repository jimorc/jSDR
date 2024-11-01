package sdr

import "github.com/pothosware/go-soapy-sdr/pkg/device"

// GetFrequencyRanges retrieves the frequency ranges supported by the device.
func (dev *StubDevice) GetFrequencyRanges(_ device.Direction, _ uint) []device.SDRRange {
	switch dev.Args["serial"] {
	case "2":
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
