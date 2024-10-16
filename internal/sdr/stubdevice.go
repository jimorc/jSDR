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

func (dev *StubDevice) ListAntennas(direction device.Direction, _ uint) []string {
	if direction == device.DirectionRX {
		return []string{"RX"}
	} else {
		return []string{}
	}
}

// GetAntennas returns the currently selected antenna for the specified direction and channel number.
// GetAntennas is misnamed as only one antenna can be selected at a time.
func (dev *StubDevice) GetAntennas(direction device.Direction, _ uint) string {
	if direction == device.DirectionRX {
		return "RX"
	} else {
		return ""
	}
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
