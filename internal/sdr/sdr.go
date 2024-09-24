package sdr

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// Sdr represents the SDR device.
type Sdr struct {
	Device      *device.SDRDevice
	SampleRates []string
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

// EnumerateWithoutAudio returns a map of SDR devices, not including any audio device.
func EnumerateWithoutAudio(log *logger.Logger) map[string]map[string]string {
	var sdrs map[string]map[string]string = make(map[string]map[string]string, 0)

	eSdrs := device.Enumerate(nil)
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
// Construction args should be as explicit as possible (i.e. include all values retrieved by EnumberateWithoutAudio).
func Make(args map[string]string, log *logger.Logger) (*Sdr, error) {
	log.Logf(logger.Debug, "Making device with label: %s\n", args["label"])
	dev, err := device.Make(args)
	if err != nil {
		log.Logf(logger.Error, "Error encountered trying to make device: %s\n", err.Error())
		return nil, err
	}
	sdr := &Sdr{Device: dev}
	log.Logf(logger.Debug, "Made SDR with hardware key: %s\n", sdr.Device.GetHardwareKey())
	return sdr, nil
}

// GetSampleRates retrieves a string slice of sample rates based on the sample rate ranges for the SDR.
func (sdr *Sdr) GetSampleRates(log *logger.Logger) []string {
	if sdr.SampleRates == nil {
		sampleRateRanges := sdr.Device.GetSampleRateRange(device.DirectionRX, 0)
		var rMsg strings.Builder
		if len(sampleRateRanges) == 0 {
			rMsg.WriteString("There are no sample rate ranges for the specified SDR\n")
		} else {
			rMsg.WriteString("Sample Rate ranges:\n")
			for _, srR := range sampleRateRanges {
				rMsg.WriteString(fmt.Sprintf("         %v\n", srR))
			}
		}
		log.Log(logger.Debug, rMsg.String())
		sdr.SampleRates = getSampleRatesAsStrings(sampleRateRanges, log)
	}
	return sdr.SampleRates
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
