package sdrdevice_test

import (
	"testing"

	"fyne.io/fyne/v2/app"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdrdevice"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromApp(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	app.NewWithID("com.github.jimorc.jsdrtest")
	var sdrDevice sdrdevice.SdrDevice
	var sdrDev2 sdrdevice.SdrDevice
	sdrDevice.Device = "Generic RTL2832U OEM :: 00000101"
	//	sdrDevice.SaveToApp(testLogger)

	sdrDev2.LoadFromApp(testLogger)
	assert.Equal(t, "Generic RTL2832U OEM :: 00000101", sdrDev2.Device)
	assert.Equal(t, sdrDevice.Device, sdrDev2.Device)
}

func TestClear(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	app.NewWithID("com.github.jimorc.jsdrtest")
	var sdrDevice sdrdevice.SdrDevice
	sdrDevice.Clear(testLogger)
	assert.Equal(t, "", sdrDevice.Device)
}
