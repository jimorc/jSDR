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
	logChannelInfo(sdr, direction, channel, log)
	logAntennaInfo(sdr, direction, channel, log)
	exerciseChannelBandwidth(sdr, direction, channel, log)
	exerciseGain(sdr, direction, channel, log)
	exerciseSampleRate(sdr, direction, channel, log)
	exerciseFrequencies(sdr, direction, channel, log)
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

func logChannelInfo(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	channelInfo := sdr.GetChannelInfo(direction, channel)
	if len(channelInfo) == 0 {
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Channel#%d ChannelInfo: none\n", channel))
	} else {
		var infoMsg strings.Builder
		infoMsg.WriteString(fmt.Sprintf("Channel#%d ChannelInfo:\n", channel))
		for k, v := range channelInfo {
			infoMsg.WriteString(fmt.Sprintf("         %s: %s\n", k, v))
		}
	}
}

func logAntennaInfo(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	antennas := sdr.ListAntennas(direction, channel)
	if len(antennas) == 0 {
		log.Log(logger.NewLogMessage(logger.Info, "Antennas: none\n"))
	} else {
		var aMsg strings.Builder
		aMsg.WriteString("Antennas:\n")
		for i, antenna := range antennas {
			aMsg.WriteString(fmt.Sprintf("         Antenna#%d: %s\n", i, antenna))
		}
		log.Log(logger.NewLogMessage(logger.Info, aMsg.String()))
	}
}

func exerciseChannelBandwidth(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	log.Log(logger.NewLogMessageWithFormat(logger.Info,
		"Channel#%d Baseband filter width: %.0f Hz\n", channel, sdr.GetBandwidth(direction, channel)))

	bandwidthRanges := sdr.GetBandwidthRanges(direction, channel)
	if len(bandwidthRanges) == 0 {
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Channel#%d Bandwidth Ranges: none\n", channel))
	} else {
		var bMsg strings.Builder
		bMsg.WriteString(fmt.Sprintf("Channel#%d Bandwidth Ranges:\n", channel))
		for i, bRange := range bandwidthRanges {
			bMsg.WriteString(fmt.Sprintf("         Bandwidth Range#%d: %v\n", i, bRange))
		}
		log.Log(logger.NewLogMessage(logger.Info, bMsg.String()))

		log.Log(logger.NewLogMessage(logger.Info, "Setting bandwidth to one half first range\n"))
		err := sdr.SetBandwidth(direction, channel, bandwidthRanges[0].Maximum/2.0)
		if err != nil {
			log.Log(logger.NewLogMessageWithFormat(logger.Error, "Error encountered while trying to set bandwidth: %sn",
				err.Error()))
		}
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Bandwidth is now %.0f\n", sdr.GetBandwidth(direction, channel)))
	}
}

func exerciseGain(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	hasAutoGainMode := sdr.HasGainMode(direction, channel)
	log.Log(logger.NewLogMessageWithFormat(logger.Info, "Channel#%d HasGainMode (Automatic gain possible): %v\n",
		channel, hasAutoGainMode))
	if hasAutoGainMode {
		autoGainEnabled := sdr.GetGainMode(direction, channel)
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Channel#%d Automatic Gain Enabled: %v\n",
			channel, autoGainEnabled))
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Current gain = %f\n",
			sdr.GetGain(direction, channel)))
		log.Log(logger.NewLogMessage(logger.Info, "Toggling auto gain\n"))
		err := sdr.SetGainMode(direction, channel, !autoGainEnabled)
		if err != nil {
			log.Log(logger.NewLogMessageWithFormat(logger.Error, "Error in call to SetGainMode: %s\n", err.Error()))
			return
		}
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Channel#%d Automatic Gain Enabled now: %v\n",
			channel, sdr.GetGainMode(direction, channel)))
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Current gain = %f\n",
			sdr.GetGain(direction, channel)))
	}
	gains := sdr.ListGains(direction, channel)
	log.Log(logger.NewLogMessageWithFormat(logger.Info, "Number of gain elements: %d\n", len(gains)))
	if len(gains) > 0 {
		var gMsg strings.Builder
		gMsg.WriteString(fmt.Sprintf("Channel#%d Gain Elements:\n", channel))
		for _, gain := range gains {
			gMsg.WriteString(fmt.Sprintf("         Element: %s\n", gain))
			gMsg.WriteString(fmt.Sprintf("             Range: %v\n", sdr.GetGainElementRange(direction, channel, gain)))
		}
		log.Log(logger.NewLogMessage(logger.Info, gMsg.String()))

		err := sdr.SetGainMode(direction, channel, false)
		if err != nil {
			log.Log(logger.NewLogMessageWithFormat(logger.Error, "Error setting auto gain off: %s\n", err.Error()))
			return
		} else {
			log.Log(logger.NewLogMessage(logger.Info, "Have set auto gain off\n"))
		}
		for _, gain := range gains {
			log.Log(logger.NewLogMessageWithFormat(logger.Info, "Setting gain for element: %s to 20 db\n", gain))
			err := sdr.SetGainElement(direction, channel, gain, 20.0)
			if err != nil {
				log.Log(logger.NewLogMessageWithFormat(logger.Error, "Error when setting gain for element: %s: %s",
					gain, err.Error()))
				return
			}
			eltGain := sdr.GetGainElement(direction, channel, gain)
			log.Log(logger.NewLogMessageWithFormat(logger.Info, "Gain for element %s is set to %.0f db\n", gain, eltGain))
			err = sdr.SetGainElement(direction, channel, gain, 0.0)
			log.Log(logger.NewLogMessageWithFormat(logger.Info, "Have reset gain for element: %s to 0 db\n", gain))
		}
		log.Log(logger.NewLogMessage(logger.Info, "Setting overall gain to 25 db\n"))
		err = sdr.SetGain(direction, channel, 25.0)
		if err != nil {
			log.Log(logger.NewLogMessageWithFormat(logger.Error, "Error when setting gain: %s", err.Error()))
			return
		}
		var gainMsg strings.Builder
		gainMsg.WriteString(fmt.Sprintf("Overall gain set to: %.0f db\n", sdr.GetGain(direction, channel)))
		for _, gain := range gains {
			eltGain := sdr.GetGainElement(direction, channel, gain)
			gainMsg.WriteString(fmt.Sprintf("         %s gain is %.0f db\n", gain, eltGain))
		}
		log.Log(logger.NewLogMessage(logger.Info, gainMsg.String()))
	}
}

func exerciseSampleRate(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	sampleRanges := sdr.GetSampleRateRange(direction, channel)
	if len(sampleRanges) == 0 {
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Channel#%d has no sample rate ranges\n", channel))
	} else {
		var sMsg strings.Builder
		sMsg.WriteString(fmt.Sprintf("Sample Rate Ranges for Channel#%d:\n", channel))
		for _, rng := range sampleRanges {
			sMsg.WriteString(fmt.Sprintf("         %v\n", rng))
		}
		log.Log(logger.NewLogMessage(logger.Info, sMsg.String()))
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Channel#%d Sample Rate: %.0f\n", channel,
			sdr.GetSampleRate(direction, channel)))
	}
	log.Log(logger.NewLogMessage(logger.Info, "Setting sample rate to 1024000\n"))
	err := sdr.SetSampleRate(direction, channel, 1024000.0)
	if err != nil {
		log.Log(logger.NewLogMessageWithFormat(logger.Error, "Error while setting sample rate: %s", err.Error()))
	}
	log.Log(logger.NewLogMessageWithFormat(logger.Info, "Sample Rate is now %.0f\n",
		sdr.GetSampleRate(direction, channel)))
}

func exerciseFrequencies(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	args := sdr.GetFrequencyArgsInfo(direction, channel)
	var aMsg strings.Builder
	if len(args) == 0 {
		log.Log(logger.NewLogMessage(logger.Info, "Frequency Args Info: none\n"))
	} else {
		aMsg.WriteString("Frequency Args Info:\n")
		for _, arg := range args {
			aMsg.WriteString(fmt.Sprintf("         %v\n", arg))
		}
		log.Log(logger.NewLogMessage(logger.Info, aMsg.String()))
	}
	freqRanges := sdr.GetFrequencyRange(direction, channel)
	if len(freqRanges) == 0 {
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Channel#%d has no frequency ranges\n", channel))
		return
	}
	var fMsg strings.Builder
	fMsg.WriteString(fmt.Sprintf("Channel#%d Frequency Ranges:\n", channel))
	for _, fRange := range freqRanges {
		fMsg.WriteString(fmt.Sprintf("         %v\n", fRange))
	}
	log.Log(logger.NewLogMessage(logger.Info, fMsg.String()))

	tuneableElts := sdr.ListFrequencies(direction, channel)
	if len(tuneableElts) == 0 {
		log.Log(logger.NewLogMessageWithFormat(logger.Info, "Channel#%d has no tuneable elements\n", channel))
	} else {
		var tMsg strings.Builder
		tMsg.WriteString(fmt.Sprintf("Tuneable elements for channel#%d:\n", channel))
		for _, elt := range tuneableElts {
			tMsg.WriteString(fmt.Sprintf("         %s\n", elt))
			comp := sdr.GetFrequencyComponent(direction, channel, elt)
			if elt == "CORR" {
				tMsg.WriteString(fmt.Sprintf("             Correction: %.0f PPM\n", comp))
				tMsg.WriteString("             Setting CORR to 50 PPM\n")
				err := sdr.SetFrequencyComponent(direction, channel, "CORR", 50.0, map[string]string{})
				if err != nil {
					log.Log(logger.NewLogMessageWithFormat(logger.Error,
						"Error encountered setting CORR: %s", err.Error()))
					return
				}
				tMsg.WriteString(fmt.Sprintf("             CORR now: %.0f PPM\n",
					sdr.GetFrequencyComponent(direction, channel, "CORR")))
			} else {
				tMsg.WriteString(fmt.Sprintf("             Frequency: %.0f Hz\n", comp))
				tMsg.WriteString("             Center Frequency set to 50 MHz\n")
				err := sdr.SetFrequencyComponent(direction, channel, elt, 50000000.0, map[string]string{})
				if err != nil {
					log.Log(logger.NewLogMessageWithFormat(logger.Error,
						"Error encountered setting component %d frequency: %s\n",
						elt, err.Error()))
					return
				}
				tMsg.WriteString(fmt.Sprintf("             Center Frequency now: %.0f Hz\n",
					sdr.GetFrequencyComponent(direction, channel, elt)))
			}
			rngs := sdr.GetFrequencyRangeComponent(direction, channel, elt)
			for _, rng := range rngs {
				tMsg.WriteString(fmt.Sprintf("             Range: %v\n", rng))
			}
		}
		log.Log(logger.NewLogMessage(logger.Info, tMsg.String()))
	}
}
