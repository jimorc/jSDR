package sdr

import (
	"errors"
	"fmt"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// agcEnabled stores current AGC enabled value.
var agcEnabled bool

// overallGain stores the current overall gain value.
var maxOverallGain float64 = 50.

type eltInfo struct {
	gain      float64
	gainRange device.SDRRange
}

var eltGains = map[string]eltInfo{
	"RF": {gain: maxOverallGain / 2, gainRange: device.SDRRange{Minimum: 0, Maximum: maxOverallGain / 2, Step: 1}},
	"IF": {gain: maxOverallGain / 2, gainRange: device.SDRRange{Minimum: 0, Maximum: maxOverallGain / 2, Step: 0}},
}

// SupportsAGC returns whether the device supports AGC or not.
//
// Returns true if device supports automatic gain control.
func (dev *StubDevice) SupportsAGC(direction device.Direction, _ uint) bool {
	if direction == device.DirectionRX {
		return true
	}
	return false
}

// AgcIsEnabled returns whether AGC is currently enabled or not.
//
// Value returned is based on the device's serial number to allow testing of sdr.AgcIsEnabled.
func (dev *StubDevice) AgcIsEnabled(direction device.Direction, _ uint) bool {
	switch dev.Args["serial"] {
	case "1":
		return true
	case "2":
		return agcEnabled
	default:
		return false
	}
}

func (dev *StubDevice) EnableAgc(_ device.Direction, _ uint, enable bool) error {
	switch dev.Args["serial"] {
	case "1":
		return errors.New("could not enable Agc")
	case "2":
		agcEnabled = enable
		return nil
	}
	return errors.New("invalid serial number for StubDevice.EnableAgc")
}

// GetGainElementNames returns a list of names for the gain elements for the specified direction and channel.
func (dev *StubDevice) GetGainElementNames(_ device.Direction, _ uint) []string {
	var keys []string
	for k := range eltGains {
		keys = append(keys, k)
	}
	return keys
}

// GetOverallGain returns the overall gain for the specified direction and channel.
func (dev *StubDevice) GetOverallGain(_ device.Direction, _ uint) float64 {
	switch dev.Args["serial"] {
	case "1":
		return 50.0
	default:
		gain := 0.0
		for _, eltGain := range eltGains {
			gain += eltGain.gain
		}
		return gain
	}
}

// SetOverallGain sets the overall gain for the specified direction and channel.
//
// The overall gain is distributed automatically across the available elements.
func (dev *StubDevice) SetOverallGain(_ device.Direction, _ uint, overallGain float64) error {
	if overallGain < 0. || overallGain > maxOverallGain {
		return errors.New(fmt.Sprintf("Requested overall gain = %.1f dB, but must be between 0.0 and %.1f dB.", overallGain, maxOverallGain))
	}
	numElts := len(eltGains)
	for k := range eltGains {
		eltGains[k] = eltInfo{gain: overallGain / float64(numElts), gainRange: eltGains[k].gainRange}
	}
	return nil
}

// GetElementGain returns the gain in dB for the specified element name.
//
// If an error occurs, such as the element name does not exist, then 0.0 and an error are returned. Note that 0.0 may also be a valid value,
// so don't just check the gain value.
func (dev *StubDevice) GetElementGain(_ device.Direction, _ uint, eltName string) (float64, error) {
	switch dev.Args["serial"] {
	case "2":
		if gain, ok := eltGains[eltName]; ok {
			return gain.gain, nil
		} else {
			return 0.0, errors.New(fmt.Sprintf("gain element '%s' is invalid", eltName))
		}
	default:
		return 25.0, nil
	}
}

// SetElementGain sets the gain for the specified element to the specified value.
//
// Returns an error if the requested gain is outside the allowable range.
func (dev *StubDevice) SetElementGain(direction device.Direction, channel uint, eltName string, gain float64) error {
	gainRange := dev.GetElementGainRange(direction, channel, eltName)
	if gain < gainRange.Minimum || gain > gainRange.Maximum {
		return errors.New(fmt.Sprintf("cannot set gain for element: %s to %.1f. Requested gain is outside the allowable range: %.1f to %.1f",
			eltName, gain, gainRange.Minimum, gainRange.Maximum))
	}
	eltGains[eltName] = eltInfo{gain, eltGains[eltName].gainRange}
	return nil
}

// GetElementGainRange returns the SDRRange for the specified gain element.
func (dev *StubDevice) GetElementGainRange(_ device.Direction, _ uint, eltName string) device.SDRRange {
	switch dev.Args["serial"] {
	case "2":
		gain, _ := eltGains[eltName]
		return gain.gainRange
	}
	return device.SDRRange{Minimum: 0, Maximum: 0, Step: 0}
}
