package main

import (
	"fmt"
	"os"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

func main() {
	logFile := "enumerate_sdrs.log"
	log, err := logger.NewFileLogger(logFile)
	if err != nil {
		fmt.Printf("Error trying to open log file '%s': %s\n", logFile, err.Error())
		os.Exit(1)
	}
	defer log.Close()

	devices := device.Enumerate(nil)
	if len(devices) == 0 {
		log.Log(logger.NewLogMessage(logger.Error, "There are no attached SDRs\n"))
	} else {
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "There are %d SDRs attached\n", len(devices)))
	}
}
