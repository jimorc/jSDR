package sdr

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

func (sD *SoapyDevice) GetAntennaNames(direction device.Direction, channel uint) []string {
	return sD.Device.Device.ListAntennas(direction, channel)
}

// GetCurrentAntenna returns the currently selected antenna for the specified direction and channel number.
func (sD *SoapyDevice) GetCurrentAntenna(direction device.Direction, channel uint) string {
	return sD.Device.Device.GetAntennas(direction, channel)
}
