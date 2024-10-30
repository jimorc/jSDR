// Package sdr provides interfaces and functions that allow multiple device types.
//
// The initial device types are:
//
//	SoapyDevice for SoapySDR devices.
//
//	StubDevice for testing of the various sdr functions.
package sdr

import (
	"fmt"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// Enumerate interface specifies the function for enumerating attached devices.
type Enumerate interface {
	Enumerate(args map[string]string) []map[string]string
}

// MakeDevice interface specifies the methods for creating and destroying an SDR device.
type MakeDevice interface {
	Make(args map[string]string) error
	Unmake() error
	GetHardwareKey() string
}

// KeyValues interface specifies the methods for retrieving SDR information.
type KeyValues interface {
	GetHardwareKey() string
}

// Sdr represents the SDR device.
type Sdr struct {
	Device           *device.SDRDevice
	DeviceProperties map[string]string
	SampleRates      []string
	SampleRate       float64
	Antennas         []string
	Antenna          string
}

// EnumerateWithoutAudio returns a map of SDR devices, not including any audio device.
func EnumerateWithoutAudio(sdrD Enumerate, log *logger.Logger) map[string]map[string]string {
	var sdrs map[string]map[string]string = make(map[string]map[string]string, 0)

	eSdrs := sdrD.Enumerate(nil)
	for _, dev := range eSdrs {
		if dev["driver"] != "audio" {
			sdrs[dev["label"]] = dev
		}
	}
	var sMsg strings.Builder
	if len(sdrs) == 0 {
		sMsg.WriteString("Attached SDRs: none\n")
	} else {
		sMsg.WriteString("Attached SDRs:\n")
		for k := range sdrs {
			sMsg.WriteString(fmt.Sprintf("         %s\n", k))
		}
	}
	log.Log(logger.Debug, sMsg.String())
	return sdrs
}

// Make makes a new device given construction args.
//
// Construction args should be as explicit as possible (i.e. include all values retrieved by
// EnumerateWithoutAudio). args should contain a label value.
func Make(sdrD MakeDevice, args map[string]string, log *logger.Logger) error {
	log.Logf(logger.Debug, "Making device with label: %s\n", args["label"])
	err := sdrD.Make(args)
	if err != nil {
		log.Logf(logger.Error, "Error encountered trying to make device: %s\n", err.Error())
		return err
	}
	log.Logf(logger.Debug, "Made SDR with hardware key: %s\n", sdrD.GetHardwareKey())
	return nil
}

// Unmake frees up any assets associated with the SDR device.
//
// No sdr calls should be made after Unmake is called.
func Unmake(sdrD MakeDevice, log *logger.Logger) error {
	log.Log(logger.Debug, "Attempting to unmake an SDR device\n")
	err := sdrD.Unmake()
	if err != nil {
		log.Logf(logger.Error, "Error attempting to unmake an SDR device: %s\n", err.Error())
	}
	return err
}

// GetHardwareKey returns the hardware key for the SDR device.
func (sdr *Sdr) GetHardwareKey(sdrD KeyValues) string {
	return sdrD.GetHardwareKey()
}
