package sdr

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// Frequency interface specifies the frequency methods that SDR devices must satisfy.
type Frequency interface {
	GetFrequencyRanges(device.Direction, uint) []device.SDRRange
	GetOverallCenterFrequency(device.Direction, uint) float64
	SetOverallCenterFrequency(device.Direction, uint, float64, map[string]string) error
	GetTunableElementNames(device.Direction, uint) []string
	GetTunableElementFrequencyRanges(device.Direction, uint, string) []device.SDRRange
	GetTunableElementFrequency(device.Direction, uint, string) float64
	SetTunableElementFrequency(device.Direction, uint, string, float64) error
}

// GetFrequencyRanges retrieves the slice of frequency ranges that the specified devices supports.
//
// There may be one or more ranges depending on the SDR device.
//
// Returns the frequency ranges for RX channel 0, or an error if there are no frequency ranges.
func GetFrequencyRanges(sdrD Frequency, log *logger.Logger) ([]device.SDRRange, error) {
	frequencyRanges := sdrD.GetFrequencyRanges(device.DirectionRX, 0)
	if len(frequencyRanges) == 0 {
		log.Log(logger.Error, "The attached SDR seems defective; there are no specified frequency ranges.\n")
		return frequencyRanges, errors.New("the attached SDR seems defective; there are no specified frequency ranges")
	}
	var frequenciesStr strings.Builder
	frequenciesStr.WriteString("Frequency Ranges:\n")
	for _, freqRange := range frequencyRanges {
		frequenciesStr.WriteString(fmt.Sprintf("         %v\n", freqRange))
	}
	log.Log(logger.Debug, frequenciesStr.String())
	return frequencyRanges, nil
}

// GetTunableElementNames retrieves the list of tunable elements by name for the device.
//
// These elements will be in the order from RF to baseband.
func GetTunableElementNames(sdrD Frequency, log *logger.Logger) []string {
	elts := sdrD.GetTunableElementNames(device.DirectionRX, 0)
	if len(elts) == 0 {
		log.Log(logger.Debug, "Device has no tunable frequency elements.\n")
	} else {
		var tMsg strings.Builder
		tMsg.WriteString("Tunable elements:\n")
		for _, elt := range elts {
			tMsg.WriteString(fmt.Sprintf("         %s\n", elt))
		}
	}
	return elts
}

// GetTunableElementFrequencyRanges retrieves the freequency ranges for the specified tunable element.
//
// Ranges are retrieved for RX channel 0 only.
// If the requested tunable element name does not match a name returned by GetTunableElementNames,
// then an error is returned.
func GetTunableElementFrequencyRanges(sdrD Frequency, log *logger.Logger, tunableElement string) ([]device.SDRRange, error) {
	tElts := sdrD.GetTunableElementNames(device.DirectionRX, 0)
	if !slices.Contains(tElts, tunableElement) {
		var eMsg strings.Builder
		eMsg.WriteString(fmt.Sprintf("Invalid tunable element name: %s\n", tunableElement))
		eMsg.WriteString(fmt.Sprintf("Tunable element names are: %v\n", tElts))
		log.Logf(logger.Error, fmt.Sprintf("Invalid "))
		return []device.SDRRange{}, fmt.Errorf("invalid tunable element name: %s", tunableElement)
	}
	fRanges := sdrD.GetTunableElementFrequencyRanges(device.DirectionRX, 0, tunableElement)
	var rMsg strings.Builder
	rMsg.WriteString(fmt.Sprintf("FrequencyRanges for tunable element: %s\n", tunableElement))
	for _, fR := range fRanges {
		rMsg.WriteString(fmt.Sprintf("         %v\n", fR))
	}
	log.Log(logger.Debug, rMsg.String())
	return fRanges, nil
}

// GetTunableElementFrequency retrieves the tuned frequency in Hz for the named tunable element.
//
// Frequency is retrieved for RX channel 0 only.
// If the requested tunable element name does not match a name returned by GetTunableElementNames,
// then an error is returned.
func GetTunableElementFrequency(sdrD Frequency, log *logger.Logger, name string) (float64, error) {
	tElts := sdrD.GetTunableElementNames(device.DirectionRX, 0)
	if !slices.Contains(tElts, name) {
		var eMsg strings.Builder
		eMsg.WriteString(fmt.Sprintf("Invalid tunable element name: %s\n", name))
		eMsg.WriteString(fmt.Sprintf("Tunable element names are: %v\n", tElts))
		log.Logf(logger.Error, fmt.Sprintf("Invalid "))
		return 0.0, fmt.Errorf("invalid tunable element name: %s", name)
	}
	eltFreq := sdrD.GetTunableElementFrequency(device.DirectionRX, 0, name)
	log.Logf(logger.Debug, fmt.Sprintf("Current frequency for element %s: %.1f\n", name, eltFreq))
	return eltFreq, nil
}

// SetTunableElementFrequency sets the frequency for the named tunable element to the specified value.
//
// Returns an error the element name is invalid or the requested frequency is not within the element's
// tunable range.
// Returns nil on success.
func SetTunableElementFrequency(sdrD Frequency, log *logger.Logger, name string, freq float64) error {
	tElts := sdrD.GetTunableElementNames(device.DirectionRX, 0)
	if !slices.Contains(tElts, name) {
		var eMsg strings.Builder
		eMsg.WriteString(fmt.Sprintf("Invalid tunable element name: %s\n", name))
		eMsg.WriteString(fmt.Sprintf("Tunable element names are: %v\n", tElts))
		log.Logf(logger.Error, eMsg.String())
		return fmt.Errorf("cannot set frequency. Invalid tunable element name: %s", name)
	}

	eltRanges, err := GetTunableElementFrequencyRanges(sdrD, log, name)
	if err != nil {
		return errors.New("cannot set frequency. Cannot retrieve tunable element frequency ranges")
	}
	if !withinRanges(eltRanges, freq) {
		return errors.New("cannot set frequency. Requested frequency not within element's frequency ranges")
	}

	sdrD.SetTunableElementFrequency(device.DirectionRX, 0, name, freq)
	log.Logf(logger.Debug, fmt.Sprintf("Setting element %s frequency to %.1f\n", name, freq))
	sdrD.GetTunableElementFrequency(device.DirectionRX, 0, name)
	return nil
}

// GetOverallCenterFrequency retrieves the overall center frequency in Hz for RX channel 0.
func GetOverallCenterFrequency(sdrD Frequency, log *logger.Logger) float64 {
	currentFreq := sdrD.GetOverallCenterFrequency(device.DirectionRX, 0)
	log.Logf(logger.Debug, "Center frequency: %.1f\n", currentFreq)
	return currentFreq
}

// SetOverallCenterFrequency sets the overall center frequency for RX channel 0.
func SetOverallCenterFrequency(sdrD Frequency, log *logger.Logger, newFreq float64, args map[string]string) error {
	freqRanges, err := GetFrequencyRanges(sdrD, log)
	if err != nil {
		log.Logf(logger.Error, "There are no frequency ranges for this device.\n")
		return fmt.Errorf("cannot set overall center frequency to %.1f.\nThere are no frequency ranges for this device", newFreq)
	}
	if !withinRanges(freqRanges, newFreq) {
		var rMsg strings.Builder
		rMsg.WriteString(fmt.Sprintf("Cannot set overall center frequency to %.1f\n"+
			"Requested frequency is not within the overall frequency ranges:\n", newFreq))
		for _, fR := range freqRanges {
			rMsg.WriteString(fmt.Sprintf("         %v\n", fR))
		}
		log.Log(logger.Error, rMsg.String())
		return fmt.Errorf("requested frequency: %.1f is not within the frequency ranges for this device", newFreq)
	}
	err = sdrD.SetOverallCenterFrequency(device.DirectionRX, 0, newFreq, args)
	if err != nil {
		log.Logf(logger.Error, "Cannot set requested overall center frequency: %.1f: %s\n", newFreq, err.Error())
		return fmt.Errorf("cannot set requested overall center frequency: %.1f: %s", newFreq, err.Error())
	}
	return nil
}

func withinRanges(ranges []device.SDRRange, value float64) bool {
	for _, r := range ranges {
		if value >= r.Minimum && value <= r.Maximum {
			return true
		}
	}
	return false
}
