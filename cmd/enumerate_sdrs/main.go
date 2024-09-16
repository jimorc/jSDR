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
	}
}
