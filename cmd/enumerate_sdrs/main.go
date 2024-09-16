package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

func main() {
	logFile := os.Getenv("HOME") + "/enumerate_sdrs.log"
	log, err := logger.NewFileLogger(logFile)
	if err != nil {
		fmt.Printf("Error trying to open log file '%s': %s\n", logFile, err.Error())
		os.Exit(1)
	}
	log.SetMaxLevel(logger.Debug)
	defer log.Close()

	devices := device.Enumerate(nil)
	log.Log(logger.NewLogMessageWithFormat(logger.Info, "Found %d attached SDR(s)\n", len(devices)))

	for i, dev := range devices {
		var devInfo strings.Builder
		devInfo.WriteString(fmt.Sprintf("Device %d\n", i))
		for k, v := range dev {
			devInfo.WriteString(fmt.Sprintf("         %s: %s\n", k, v))
		}
		log.Log(logger.NewLogMessage(logger.Info, devInfo.String()))

		// Open device
		log.Log(logger.NewLogMessageWithFormat(logger.Debug,
			"Making device with label: '%s'\n", dev["label"]))
		sdr, err := device.Make(dev)
		if err != nil {
			log.Log(logger.NewLogMessageWithFormat(logger.Error,
				"Unable to make device with label: %s: %s\n", dev["label"], err.Error))
		}
		if sdr == nil {
			log.Log(logger.NewLogMessage(logger.Error, "Could not make SDR\n"))
		}
		defer func() {
			err := sdr.Unmake()
			if err != nil {
				log.Log(logger.NewLogMessageWithFormat(logger.Error,
					"Could not Unmake SDR with label: %s: %s\n",
					dev["label"], err.Error()))
				fmt.Println("Unable to Unmake a device. See log file for more info.")
				os.Exit(1)
			}
			log.Log(logger.NewLogMessageWithFormat(logger.Debug,
				"Device with label: `%s` was unmade.\n",
				dev["label"]))
		}()

		displayHardwareInfo(sdr, log)
	}
}

func displayHardwareInfo(sdr *device.SDRDevice, log *logger.Logger) {
	var hwInfo strings.Builder
	hwInfo.WriteString(fmt.Sprintln("Hardware Info:"))
	hwInfo.WriteString(fmt.Sprintf("         Driver Key: %s\n", sdr.GetDriverKey()))
	hwInfo.WriteString(fmt.Sprintf("         Hardware Key: %s\n", sdr.GetHardwareKey()))
	hardwareInfo := sdr.GetHardwareInfo()

	for k, v := range hardwareInfo {
		hwInfo.WriteString(fmt.Sprintf("         %s: %s\n", k, v))
	}
	log.Log(logger.NewLogMessage(logger.Info, hwInfo.String()))
}
