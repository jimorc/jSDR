package sdr

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// StubDevice provides a stub interface for testing the soapySDR interface.
// The fields in the struct allow loading of data for test purposes to allow simple method returns.
type StubDevice struct {
	Device     *Sdr
	Devices    []map[string]string
	Args       map[string]string
	sampleRate float64
}

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

// Enumerate returns a slice of map[string]string values representing the available devices. These
// values must be preloaded into the Devices property of the StubDevice struct before Enumerate is called.
func (dev StubDevice) Enumerate(args map[string]string) []map[string]string {
	return dev.Devices
}

// Make simply sets StubDevice.Device field to a fake device pointer and returns an error or nil based on
// the args passed into Make.
//
// If args are provided, then Make sets the Device.Device field to a fake device pointer. If no args are provided,
// then an error is returned.
func (dev *StubDevice) Make(args map[string]string) error {
	if len(args) == 0 {
		return errors.New("No arguments provided")
	} else {
		dev.Args = args
		fakeDev := 127
		dev.Device = &Sdr{Device: (*device.SDRDevice)(unsafe.Pointer(&fakeDev)),
			DeviceProperties: args}
		return nil
	}
}

// Unmake returns nil if previous call to Make was successful; otherwise, returns an error.
// Since StubDevice is used for testing, an actual SDR is never created.
func (dev *StubDevice) Unmake() error {
	if dev.Device == nil {
		return errors.New("No device to unmake")
	} else {
		dev.Device = nil
		dev.Args = nil
		return nil
	}
}

// GetHardwareKey returns a string containing a fake hardware key.
func (dev StubDevice) GetHardwareKey() string {
	return "hardKey"
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
		return errors.New("Could not enable Agc")
	case "2":
		agcEnabled = enable
		return nil
	}
	return errors.New("Invalid serial number for StubDevice.EnableAgc")
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
			return 0.0, errors.New(fmt.Sprintf("Gain element '%s' is invalid", eltName))
		}
	default:
		return 25.0, nil
	}
}

// GetElementGainRange returns the SDRRange for the specified gain element.
func (dev *StubDevice) GetElementGainRange(_ device.Direction, _ uint, eltName string) device.SDRRange {
	switch dev.Args["serial"] {
	case "2":
		return device.SDRRange{Minimum: 0, Maximum: 25, Step: 0}
	}
	return device.SDRRange{Minimum: 0, Maximum: 0, Step: 0}
}
