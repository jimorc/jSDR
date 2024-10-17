package sdr

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

func (sD *SoapyDevice) ListAntennas(direction device.Direction, channel uint) []string {
	return sD.Device.Device.ListAntennas(direction, channel)
}

// GetAntennas returns the currently selected antenna for the specified direction and channel number.
// GetAntennas is misnamed as only one antenna can be selected at a time.
func (sD *SoapyDevice) GetAntennas(direction device.Direction, channel uint) string {
	return sD.Device.Device.GetAntennas(direction, channel)
}
