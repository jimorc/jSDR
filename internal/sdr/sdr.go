package sdr

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

type Enumerate interface {
	Enumerate(args map[string]string) []map[string]string
}

type MakeDevice interface {
	Make(args map[string]string) error
	Unmake() error
	GetHardwareKey() string
}

type KeyValues interface {
	GetHardwareKey() string
}

type Agc interface {
	SupportsAGC(device.Direction, uint) bool
	AgcIsEnabled(device.Direction, uint) bool
	EnableAgc(device.Direction, uint, bool) error
}

type Gain interface {
	GetGainElementNames(device.Direction, uint) []string
	GetOverallGain(device.Direction, uint) float64
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

func GetGainElementNames(sdrD Gain, log *logger.Logger) []string {
	elts := sdrD.GetGainElementNames(device.DirectionRX, 0)
	log.Logf(logger.Debug, "Gain Elements: %v\n", elts)
	return elts
}

// GetOverallGain gets the overall value of the gain elements in the chain for RX channel 0.
//
// Returns value of the gain in dB.
func GetOverallGain(sdrD Gain, log *logger.Logger) float64 {
	gain := sdrD.GetOverallGain(device.DirectionRX, 0)
	log.Logf(logger.Debug, "Overall gain: %.1f\n", gain)
	return gain
}

func SetSampleRate(sdrD SampleRates, log *logger.Logger, rate float64) error {
	GetSampleRates(sdrD, log)
	currentRate := GetSampleRate(sdrD, log)
	re, err := regexp.Compile(`[0-9]+\.[0-9]+`)
	match := re.FindString(currentRate)
	sampleRate, err := strconv.ParseFloat(match, 64)
	if err != nil {
		log.Logf(logger.Error, "Error parsing sample rate: %s\n", err.Error())
	}
	sampleRate *= 1e6
	rate = closestSampleRate(rate, log)
	if rate == sampleRate {
		log.Log(logger.Debug, "Requested rate is same as current sample rate.\n")
		return nil
	} else {
		log.Logf(logger.Debug, "Setting sample rate to %f\n", rate)
		err := sdrD.SetSampleRate(device.DirectionRX, 0, rate)
		if err != nil {
			log.Logf(logger.Error, "Error attempting to set sample rate: %s\n", err.Error())
			return err
		} else {
			match = re.FindString(GetSampleRate(sdrD, log))
			setRate, _ := strconv.ParseFloat(match, 64)
			setRate *= 1e6
			if setRate != rate {
				errMsg := fmt.Sprintf("Attempt to set sample rate to %.1f failed. Sample rate is %.1f", rate, setRate)
				log.Log(logger.Error, errMsg)
				return errors.New(errMsg)
			} else {
				log.Logf(logger.Debug, "Sample rate has been set to %f\n", setRate)
				return nil
			}
		}
	}
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
