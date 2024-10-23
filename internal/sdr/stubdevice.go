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
var overallGain float64 = 50.

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
	return []string{"RX"}
}

// GetOverallGain returns the overall gain for the specified direction and channel.
func (dev *StubDevice) GetOverallGain(_ device.Direction, _ uint) float64 {
	return overallGain
}

// SetOverallGain sets the overall gain for the specified direction and channel.
//
// The overall gain is distributed automatically across the available elements. This is currently not done in StubDevice.
func (dev *StubDevice) SetOverallGain(_ device.Direction, _ uint, gain float64) error {
	if gain < 0. || gain > 50. {
		return errors.New(fmt.Sprintf("Requested overall gain = %.1f dB, but must be between 0. and 50. dB.", gain))
	}
	overallGain = gain
	return nil
}
