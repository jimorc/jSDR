package main

import (
	"fmt"
	"os"
	"strings"
	"time"

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

		logHardwareInfo(sdr, log)
		logGPIOBanks(sdr, log)
		logSettingInfo(sdr, log)
		logUARTs(sdr, log)
		logMasterClockRate(sdr, log)
		logClockSources(sdr, log)
		logRegisters(sdr, log)
		logSensors(sdr, log)
		logTimeSources(sdr, log)
		logDirectionDetails(sdr, device.DirectionTX, log)
		logDirectionDetails(sdr, device.DirectionRX, log)
	}
}

func logHardwareInfo(sdr *device.SDRDevice, log *logger.Logger) {
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

func logGPIOBanks(sdr *device.SDRDevice, log *logger.Logger) {
	banks := sdr.ListGPIOBanks()
	if len(banks) == 0 {
		log.Log(logger.NewLogMessage(logger.Info, "GPIO Banks: none\n"))
	} else {
		var gpioBanks strings.Builder
		gpioBanks.WriteString("GPIO Banks:\n")
		for i, bank := range banks {
			gpioBanks.WriteString(fmt.Sprintf("         GPIO Bank#%d: %v\n", i, bank))
		}
		log.Log(logger.NewLogMessage(logger.Info, gpioBanks.String()))
	}
}

func logSettingInfo(sdr *device.SDRDevice, log *logger.Logger) {
	SDRSettings := sdr.GetSettingInfo()
	if len(SDRSettings) == 0 {
		log.Log(logger.NewLogMessage(logger.Info, "Settings: none"))
	} else {
		var settings strings.Builder
		for i, set := range SDRSettings {
			if i == 0 {
				settings.WriteString(fmt.Sprintf("Setting%d:\n", i))
			} else {
				settings.WriteString(fmt.Sprintf("        Setting%d:\n", i))
			}
			settings.WriteString(fmt.Sprintf("         key: %s\n", set.Key))
			settings.WriteString(fmt.Sprintf("         value: %s\n", set.Value))
			settings.WriteString(fmt.Sprintf("         name: %s\n", set.Name))
			settings.WriteString(fmt.Sprintf("         description: %s\n", set.Description))
			settings.WriteString(fmt.Sprintf("         unit: %s\n", set.Unit))
			argType := "unknown type"
			switch set.Type {
			case device.ArgInfoBool:
				argType = "bool"
			case device.ArgInfoInt:
				argType = "int"
			case device.ArgInfoFloat:
				argType = "float"
			case device.ArgInfoString:
				argType = "string"
			}
			settings.WriteString(fmt.Sprintf("         type: %s\n", argType))
			settings.WriteString(fmt.Sprintf("         range: %v\n", set.Range.ToString()))
			numOptions := set.NumOptions
			if numOptions == 0 {
				settings.WriteString(fmt.Sprintln("         options: none"))
				settings.WriteString(fmt.Sprintln("         option names: none"))
			} else {
				settings.WriteString(fmt.Sprintln("         options:"))
				for _, opt := range set.Options {
					settings.WriteString(fmt.Sprintf("            %s\n", opt))
				}
				settings.WriteString(fmt.Sprintln("            option names:"))
				for _, name := range set.OptionNames {
					settings.WriteString(fmt.Sprintf("             %s\n", name))
				}
			}
		}
		log.Log(logger.NewLogMessage(logger.Info, settings.String()))
	}
}

func logUARTs(sdr *device.SDRDevice, log *logger.Logger) {
	uarts := sdr.ListUARTs()
	if len(uarts) == 0 {
		log.Log(logger.NewLogMessage(logger.Info, "UARTs: none\n"))
	} else {
		var umsg strings.Builder
		umsg.WriteString("UARTs:\n")
		for i, uart := range uarts {
			umsg.WriteString(fmt.Sprintf("         UART#%d: %s", i, uart))
		}
		log.Log(logger.NewLogMessage(logger.Info, umsg.String()))
	}
}

func logMasterClockRate(sdr *device.SDRDevice, log *logger.Logger) {
	clockRates := sdr.GetMasterClockRates()
	if len(clockRates) == 0 {
		log.Log(logger.NewLogMessage(logger.Info, "Master Clock Rates: none\n"))
	} else {
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Master Clock Rate: %f\n", sdr.GetMasterClockRate()))
		var rMsg strings.Builder
		rMsg.WriteString("Master Clock Rates:\n")
		for _, rate := range clockRates {
			rMsg.WriteString(fmt.Sprintf("         %v\n", rate))
		}
		log.Log(logger.NewLogMessage(logger.Info, rMsg.String()))
	}
}

func logClockSources(sdr *device.SDRDevice, log *logger.Logger) {
	sources := sdr.ListClockSources()
	if len(sources) == 0 {
		log.Log(logger.NewLogMessage(logger.Info, "Clock Sources: none\n"))
	} else {
		var sMsg strings.Builder
		sMsg.WriteString("Clock Sources:\n")
		for i, source := range sources {
			sMsg.WriteString(fmt.Sprintf("         Source#%d: %s\n", i, source))
		}
		log.Log(logger.NewLogMessage(logger.Info, sMsg.String()))
	}
}

func logRegisters(sdr *device.SDRDevice, log *logger.Logger) {
	registers := sdr.ListRegisterInterfaces()
	if len(registers) == 0 {
		log.Log(logger.NewLogMessage(logger.Info, "Registers: none\n"))
	} else {
		var rMsg strings.Builder
		rMsg.WriteString("Registers:\n")
		for i, register := range registers {
			rMsg.WriteString(fmt.Sprintf("         Register#%d: %s\n", i, register))
		}
		log.Log(logger.NewLogMessage(logger.Info, rMsg.String()))
	}
}

func logSensors(sdr *device.SDRDevice, log *logger.Logger) {
	sensors := sdr.ListSensors()
	if len(sensors) == 0 {
		log.Log(logger.NewLogMessage(logger.Info, "Sensors: none\n"))
	} else {
		var sMsg strings.Builder
		sMsg.WriteString("Sensors:\n")
		for i, sensor := range sensors {
			sMsg.WriteString(fmt.Sprintf("         Sensor#%d: %s\n", i, sensor))
		}
		log.Log(logger.NewLogMessage(logger.Info, sMsg.String()))
	}
}

func logTimeSources(sdr *device.SDRDevice, log *logger.Logger) {
	sources := sdr.ListTimeSources()
	if len(sources) == 0 {
		log.Log(logger.NewLogMessage(logger.Info, "Time Sources: none\n"))
	} else {
		var tMsg strings.Builder
		tMsg.WriteString("Time Sources:\n")
		for i, source := range sources {
			tMsg.WriteString(fmt.Sprintf("         Time Source#%d: %s\n", i, source))
		}
		log.Log(logger.NewLogMessage(logger.Info, tMsg.String()))
	}

	hasHardwareTime := sdr.HasHardwareTime("")
	log.Log(logger.NewLogMessageWithFormat(logger.Info, "Has Hardware Time: %v\n", hasHardwareTime))
	if hasHardwareTime {
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Hardware Time: %d ns\n", sdr.GetHardwareTime("")))
		curTime := time.Now().UTC().Nanosecond()
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Setting Hardware Time to %d\n", curTime))
		sdr.SetHardwareTime(uint(curTime), "")
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Hardware Time Now: %d\n", sdr.GetHardwareTime("")))
		log.Log(logger.NewLogMessage(logger.Info, "Waiting 1 second\n"))
		time.Sleep(time.Second)
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Hardware Time Now: %d\n", sdr.GetHardwareTime("")))
	}
}

func logDirectionDetails(sdr *device.SDRDevice, direction device.Direction, log *logger.Logger) {
	if direction == device.DirectionTX {
		log.Log(logger.NewLogMessage(logger.Info, "Direction TX\n"))
	} else {
		log.Log(logger.NewLogMessage(logger.Info, "Direction RX\n"))
	}

	frontendMapping := sdr.GetFrontendMapping(direction)
	if len(frontendMapping) == 0 {
		log.Log(logger.NewLogMessage(logger.Info, "Frontend Mapping: none\n"))
	} else {
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Frontend Mapping: %s\n", frontendMapping))
	}

	numChannels := sdr.GetNumChannels(direction)
	log.Log(logger.NewLogMessageWithFormat(logger.Info, "Number of channels: %d\n", numChannels))
	for ch := uint(0); ch < numChannels; ch++ {
		logDirectionChannelDetails(sdr, direction, ch, log)
	}
}

func logDirectionChannelDetails(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	logChannelSettingsInfo(sdr, direction, channel, log)
}

func logChannelSettingsInfo(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	settings := sdr.GetChannelSettingInfo(direction, channel)
	if len(settings) == 0 {
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Channel#%d Settings: none\n", channel))
	} else {
		var sMsg strings.Builder
		sMsg.WriteString(fmt.Sprintf("Channel#%d Settings:\n", channel))
		for i, setting := range settings {
			sMsg.WriteString(fmt.Sprintf("         Channel#%d Setting#%d: %v\n", channel, i, setting))
		}
		log.Log(logger.NewLogMessage(logger.Info, sMsg.String()))
	}

}
