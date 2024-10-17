package sdr

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// GetAntennaNames retrieves a list of all antennas for the direction and channel
func (dev *StubDevice) GetAntennaNames(direction device.Direction, _ uint) []string {
	if direction == device.DirectionRX {
		return []string{"RX"}
	} else {
		return []string{}
	}
}

// GetCurrentAntenna returns the currently selected antenna for the specified direction and channel number.
func (dev *StubDevice) GetCurrentAntenna(direction device.Direction, _ uint) string {
	if direction == device.DirectionRX {
		return "RX"
	} else {
		return ""
	}
}
