package sdr

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// ListAntennas retrieves a list of all antennas for the direction and channel
func (dev *StubDevice) GetAntennaNames(direction device.Direction, _ uint) []string {
	if direction == device.DirectionRX {
		return []string{"RX"}
	} else {
		return []string{}
	}
}

// GetAntennas returns the currently selected antenna for the specified direction and channel number.
// GetAntennas is misnamed as only one antenna can be selected at a time.
func (dev *StubDevice) GetAntennas(direction device.Direction, _ uint) string {
	if direction == device.DirectionRX {
		return "RX"
	} else {
		return ""
	}
}
