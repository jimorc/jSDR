package sdr

import (
	"fmt"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"
)

// Sdrs is a struct containing a map of devices and their properties
type Sdrs struct {
	DevicesMap map[string]map[string]string
}

// Enumerate interface specifies the function for enumerating attached devices.
type Enumerate interface {
	Enumerate(args map[string]string) []map[string]string
}

// EnumerateSdrsWithoutAudio returns an Sdrs struct containing all SDRs connected
// to the computer except the audio device.
func EnumerateSdrsWithoutAudio(sdrD Enumerate, log *logger.Logger) Sdrs {
	var sdrs Sdrs
	sdrs.DevicesMap = make(map[string]map[string]string, 0)

	eSdrs := sdrD.Enumerate(nil)
	for _, dev := range eSdrs {
		if dev["driver"] != "audio" {
			sdrs.DevicesMap[dev["label"]] = dev
		}
	}
	var sMsg strings.Builder
	if sdrs.NumberOfSdrs() == 0 {
		sMsg.WriteString("Attached SDRs: none\n")
	} else {
		sMsg.WriteString("Attached SDRs:\n")
		for k := range sdrs.DevicesMap {
			sMsg.WriteString(fmt.Sprintf("         %s\n", k))
		}
	}
	log.Log(logger.Debug, sMsg.String())
	return sdrs
}

// Contains determines if the Sdrs struct contains the specified Sdr device.
//
// Params:
//
//	 label is the label for the SDR.
//		log is the logger to write messages to
func (s *Sdrs) Contains(label string, log *logger.Logger) bool {
	for k := range s.DevicesMap {
		if k == label {
			log.Logf(logger.Debug, "Found device with label: %s\n", label)
			return true
		}
	}
	log.Logf(logger.Debug, "Could not find device with label: %s\n", label)
	return false
}

// NumberOfSdrs returns the number of SDRs in the Sdrs struct.
func (s *Sdrs) NumberOfSdrs() int {
	return len(s.DevicesMap)
}

func (s *Sdrs) SdrLabels() []string {
	var sdrLabels []string
	for k := range s.DevicesMap {
		sdrLabels = append(sdrLabels, k)
	}
	return sdrLabels
}

func (s *Sdrs) Sdr(label string) map[string]string {
	for k, v := range s.DevicesMap {
		if k == label {
			return v
		}
	}
	return nil
}
