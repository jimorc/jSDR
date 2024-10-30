package sdr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// Frequency interface specifies the frequency methods that SDR devices must satisfy.
type Frequency interface {
	GetFrequencyRanges(device.Direction, uint) []device.SDRRange
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
