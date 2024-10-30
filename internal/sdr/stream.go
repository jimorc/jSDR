package sdr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// Stream interface specifies stream related functionality.
type Stream interface {
	GetStreamFormats(device.Direction, uint) []string
	GetNativeStreamFormat(device.Direction, uint) (string, float64)
}

// GetStreamFormats retrieves the stream formats for RX channel 0.
//
// Returns an error if there are no formats for RX channel 0.
func GetStreamFormats(sdrD Stream, log *logger.Logger) ([]string, error) {
	formats := sdrD.GetStreamFormats(device.DirectionRX, 0)
	if len(formats) == 0 {
		log.Log(logger.Error, "Channel 0 has no stream formats\n")
		err := errors.New("No stream formats retrieved for channel 0")
		return formats, err
	}
	var formatStr strings.Builder
	formatStr.WriteString("Stream Formats:\n")
	for _, format := range formats {
		formatStr.WriteString(fmt.Sprintf("         %s\n", format))
	}
	log.Log(logger.Debug, formatStr.String())
	return formats, nil
}

// GetNativeStreamFormat retrieves the native format for RX channel 0 and its full scale value.
func GetNativeStreamFormat(sdrD Stream, log *logger.Logger) (string, float64) {
	return sdrD.GetNativeStreamFormat(device.DirectionRX, 0)
}
