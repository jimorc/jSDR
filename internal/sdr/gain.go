package sdr

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// Agc interface specifies the AGC methods.
type Agc interface {
	SupportsAGC(device.Direction, uint) bool
	AgcIsEnabled(device.Direction, uint) bool
	EnableAgc(device.Direction, uint, bool) error
}

// Gain interface specifies the gain related methods.
type Gain interface {
	GetGainElementNames(device.Direction, uint) []string
	GetElementGain(device.Direction, uint, string) (float64, error)
	SetElementGain(device.Direction, uint, string, float64) error
	GetElementGainRange(device.Direction, uint, string) device.SDRRange
	GetOverallGain(device.Direction, uint) float64
	SetOverallGain(device.Direction, uint, float64) error
}

// SupportsAGC returns whether the device supports AGC or not.
//
// Returns true if device supports automatic gain control for RX channel 0.
func SupportsAGC(sdrD Agc, log *logger.Logger) bool {
	supportsAGC := sdrD.SupportsAGC(device.DirectionRX, 0)
	log.Logf(logger.Debug, "Device has gain mode: %v\n", supportsAGC)
	return supportsAGC
}

// AgcIsEnabled returns whether AGC is enabled for receive channel 0 of the device.
// You should call SupportsAGC to determine if the device supports AGC before calling AgcIsEnabled.
//
// Returns true if AGC is enabled.
func AgcIsEnabled(sdrD Agc, log *logger.Logger) bool {
	enabled := sdrD.AgcIsEnabled(device.DirectionRX, 0)
	log.Logf(logger.Debug, "AgcIsEnabled: %v\n", enabled)
	return enabled
}

// EnableAge enables or disables the device's AGC.
func EnableAgc(sdrD Agc, log *logger.Logger, enable bool) error {
	err := sdrD.EnableAgc(device.DirectionRX, 0, enable)
	if err != nil {
		log.Logf(logger.Debug, "Error returned trying to set AGC mode: %v: %s\n", enable, err.Error())
		return err
	}
	enabled := sdrD.AgcIsEnabled(device.DirectionRX, 0)
	if enable == enabled {
		log.Logf(logger.Debug, "AGC mode set to %v\n", enable)
	} else {
		log.Logf(logger.Debug, "Error attempting to set AGC mode to %v\n", enable)
	}
	return err
}

// GetGainElementNames retrieves the list of gain elements for the SDR.
func GetGainElementNames(sdrD Gain, log *logger.Logger) []string {
	elts := sdrD.GetGainElementNames(device.DirectionRX, 0)
	log.Logf(logger.Debug, "Gain Elements: %v\n", elts)
	return elts
}

// GetElementGain gets the gain for the named element in the chain for RX channel 0.
//
// Returns the gain for the specified element in dB and nil, or 0.0 and an error on error.
// Do not just check if the gain is 0.0 because this may be a valid value for the named element.
func GetElementGain(sdrD Gain, log *logger.Logger, elementName string) (float64, error) {
	gain, err := sdrD.GetElementGain(device.DirectionRX, 0, elementName)
	if err != nil {
		log.Logf(logger.Error, "Error getting gain for element: %s: %s\n", elementName, err.Error())
		return 0.0, err
	}
	log.Logf(logger.Debug, "Gain for element %s is %.1f\n", elementName, gain)
	return gain, nil
}

// GetElementGainRange retrieves the gain range for the specified device.
//
// Returns an error if the requested element does not exist.
func GetElementGainRange(sdrD Gain, log *logger.Logger, elementName string) (device.SDRRange, error) {
	elts := sdrD.GetGainElementNames(device.DirectionRX, 0)
	valid := false
	for _, elt := range elts {
		if elementName == elt {
			valid = true
			break
		}
	}
	if !valid {
		errStr := fmt.Sprintf("Gain element name: %s is invalid\n", elementName)
		log.Logf(logger.Error, errStr+"\n")
		return device.SDRRange{Minimum: 0, Maximum: 0, Step: 0}, errors.New(errStr)
	}
	return sdrD.GetElementGainRange(device.DirectionRX, 0, elementName), nil
}

// SetElementGain attempts to set the gain for the specified element to the requested value.
//
// Returns an error if:
//
//		The gain element does not exist.
//	 The requested gain is outside the gain range for the element.
func SetElementGain(sdrD Gain, log *logger.Logger, elementName string, gain float64) error {
	eltNames := sdrD.GetGainElementNames(device.DirectionRX, 0)
	if !slices.Contains(eltNames, elementName) {
		log.Logf(logger.Error, "Attempting to set gain for element: %s, but that gain element does not exist.\n"+
			"Gain elements are: %v\n", elementName, eltNames)
		return fmt.Errorf("cannot set gain for non-existent gain element: %s", elementName)
	}
	err := sdrD.SetElementGain(device.DirectionRX, 0, elementName, gain)
	if err != nil {
		log.Logf(logger.Error, fmt.Sprintf("unable to set gain for element: %s: %s\n", elementName, err))
	}
	return err
}

// GetOverallGain gets the overall value of the gain elements in the chain for RX channel 0.
//
// Returns value of the gain in dB.
func GetOverallGain(sdrD Gain, log *logger.Logger) float64 {
	gain := sdrD.GetOverallGain(device.DirectionRX, 0)
	log.Logf(logger.Debug, "Overall gain: %.1f\n", gain)
	return gain
}

// SetOverallGain sets the total overall gain of the various gain elements in the chain to the specified value in dB.
//
// Returns nil on success, or error on failure.
func SetOverallGain(sdrD Gain, log *logger.Logger, gain float64) error {
	err := sdrD.SetOverallGain(device.DirectionRX, 0, gain)
	if err != nil {
		log.Logf(logger.Error, "Could not set overall gain to %.1f: %s\n", gain, err.Error())
	} else {
		log.Logf(logger.Debug, "Have set overall gain to %.1f\n", gain)
	}
	return err
}

func closestSampleRate(sampleRate float64, log *logger.Logger) float64 {
	log.Logf(logger.Debug, "closestSampleRate called with sampleRate = %f\n", sampleRate)
	lowerRate := 0.0
	upperRate := 10000000.0
	log.Logf(logger.Debug, "lowerRate = %f, upperRate = %f\n", lowerRate, upperRate)
	for rate := range sampleRatesMap {
		if sampleRate >= rate && rate > lowerRate {
			lowerRate = rate
		}
		if sampleRate <= rate && rate < upperRate {
			upperRate = rate
		}
		log.Logf(logger.Debug, "For rate = %f, lowerRate = %f, upperRate = %f\n", rate, lowerRate, upperRate)
	}
	lowerDiff := sampleRate - lowerRate
	upperDiff := upperRate - sampleRate
	if lowerDiff < upperDiff {
		log.Logf(logger.Debug, "closestSampleRate returning %f\n", lowerRate)
		return lowerRate
	} else {
		log.Logf(logger.Debug, "closestSampleRate returning %f\n", upperRate)
		return upperRate
	}
}

func getSampleRatesAsStrings(samplewidths []device.SDRRange, log *logger.Logger) []string {
	var rates []string
	for k, v := range sampleRatesMap {
		for _, bw := range samplewidths {
			if k >= bw.Minimum && k <= bw.Maximum {
				rates = append(rates, v)
				break
			}
		}
	}
	var rMsg strings.Builder
	if len(rates) == 0 {
		rMsg.WriteString("Sample Rates found: none\n")
	} else {
		sort.Strings(rates)
		rMsg.WriteString("Sample Rates found:\n")
		for _, rate := range rates {
			rMsg.WriteString(fmt.Sprintf("         %s\n", rate))
		}
	}
	log.Log(logger.Debug, rMsg.String())
	return rates
}
