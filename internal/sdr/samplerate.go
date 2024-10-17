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

// map of possible sample rates
var sampleRatesMap = map[float64]string{
	256000.0:   "0.256 MS/s",
	512000.0:   "0.512 MS/s",
	1024000.0:  "1.024 MS/s",
	1600000.0:  "1.6 MS/s",
	2048000.0:  "2.048 MS/s",
	2400000.0:  "2.4 MS/s",
	2800000.0:  "2.8 MS/s",
	3200000.0:  "3.2 MS/s",
	4000000.0:  "4.0 MS/s",
	5000000.0:  "5.0 MS/s",
	6000000.0:  "6.0 MS/s",
	7000000.0:  "7.0 MS/s",
	8000000.0:  "8.0 MS/s",
	9000000.0:  "9.0 MS/s",
	10000000.0: "10.0 MS/s",
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
