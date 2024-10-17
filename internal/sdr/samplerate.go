package sdr

import (
	"fmt"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

type SampleRates interface {
	GetSampleRateRange(device.Direction, uint) []device.SDRRange
	GetSampleRate(device.Direction, uint) float64
	SetSampleRate(device.Direction, uint, float64) error
}

// GetSampleRate returns the sample rate that most closely matches the current sample rate for the SDR.
func GetSampleRate(sdrD SampleRates, log *logger.Logger) string {
	sampleRate := sdrD.GetSampleRate(device.DirectionRX, 0)
	log.Logf(logger.Debug, "Current sample rate: %f\n", sampleRate)
	closestRate := closestSampleRate(sampleRate, log)
	log.Logf(logger.Debug, "Closest sample rate: %f\n", closestRate)
	log.Logf(logger.Debug, "GetSampleRate returning %s\n", sampleRatesMap[closestRate])
	return sampleRatesMap[closestRate]
}

// GetSampleRates retrieves a string slice of sample rates based on the sample rate ranges for the SDR.
func GetSampleRates(sdrD SampleRates, log *logger.Logger) []string {
	sampleRateRanges := sdrD.GetSampleRateRange(device.DirectionRX, 0)
	var rMsg strings.Builder
	if len(sampleRateRanges) == 0 {
		rMsg.WriteString("There are no sample rate ranges for the specified SDR\n")
	} else {
		rMsg.WriteString("Sample Rate ranges:\n")
		for _, srR := range sampleRateRanges {
			rMsg.WriteString(fmt.Sprintf("         %v\n", srR))
		}
		log.Log(logger.Debug, rMsg.String())
	}
	return getSampleRatesAsStrings(sampleRateRanges, log)
}
