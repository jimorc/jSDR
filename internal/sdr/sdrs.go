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
	if len(sdrs.DevicesMap) == 0 {
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
