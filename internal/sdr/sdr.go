package sdr

import (
	"fmt"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

type Sdr struct {
	Device *device.SDRDevice
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
