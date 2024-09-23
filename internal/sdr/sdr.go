package sdr

import (
	"fmt"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// EnumerateWithoutAudio returns a slice of SDR device labels, not including any audio device.
func EnumerateWithoutAudio(log *logger.Logger) []string {
	var sdrLabels []string
	var sdrs []map[string]string

	eSdrs := device.Enumerate(nil)
	for _, dev := range eSdrs {
		if dev["driver"] != "audio" {
			sdrs = append(sdrs, dev)
		}
	}
	var sMsg strings.Builder
	if len(sdrs) == 0 {
		sMsg.WriteString("Attached SDRs: none\n")
	} else {
		sMsg.WriteString("Attached SDRs:\n")
		for _, dev := range sdrs {
			sdrLabels = append(sdrLabels, dev["label"])
			sMsg.WriteString(fmt.Sprintf("         %s\n", dev["label"]))
		}
	}
	log.Log(logger.Debug, sMsg.String())
	return sdrLabels
}
