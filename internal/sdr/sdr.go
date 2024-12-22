// Package sdr provides interfaces and functions that allow multiple device types.
//
// The initial device types are:
//
//	SoapyDevice for SoapySDR devices.
//
//	StubDevice for testing of the various sdr functions.
//
// Many of the function and method names are changed from those provided in go-soapy-sdr.go.
// I find many of the function and method names to be confusing in go-soapy-sdr.go For example:
// device.SetAntennas sets a single antenna on a device, not multiple antennas.
//
// These name changes are an attempt to clarify what the functions and methods do,
package sdr

import (
	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

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

var SoapyDev = &SoapyDevice{}

// Make makes a new device given construction args.
//
// Construction args should be as explicit as possible (i.e. include all values retrieved by
// EnumerateSdrsWithoutAudio). args should contain a label value.
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
