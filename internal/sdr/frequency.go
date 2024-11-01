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
	GetTunableElements(device.Direction, uint) []string
	GetTunableElementFrequencyRanges(device.Direction, uint, string) []device.SDRRange
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
		return frequencyRanges, errors.New("The attached SDR seems defective; there are no specified frequency ranges.")
	}
	var frequenciesStr strings.Builder
	frequenciesStr.WriteString("Frequency Ranges:\n")
	for _, freqRange := range frequencyRanges {
		frequenciesStr.WriteString(fmt.Sprintf("         %v\n", freqRange))
	}
	log.Log(logger.Debug, frequenciesStr.String())
	return frequencyRanges, nil
}

// GetTunableElements retrieves the list of tunable elements by name for the device.
//
// These elements will be in the order from RF to baseband.
func GetTunableElements(sdrD Frequency, log *logger.Logger) []string {
	elts := sdrD.GetTunableElements(device.DirectionRX, 0)
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
func GetTunableElementFrequencyRanges(sdrD Frequency, log *logger.Logger, tunableElement string) ([]device.SDRRange, error) {
	tElts := sdrD.GetTunableElements(device.DirectionRX, 0)
	if !slices.Contains(tElts, tunableElement) {
		var eMsg strings.Builder
		eMsg.WriteString(fmt.Sprintf("Invalid tunable element name: %s\n", tunableElement))
		eMsg.WriteString(fmt.Sprintf("Tunable element names are: %v\n", tElts))
		log.Logf(logger.Error, fmt.Sprintf("Invalid "))
		return []device.SDRRange{}, errors.New(fmt.Sprintf("Invalid tunable element name: %s", tunableElement))
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
