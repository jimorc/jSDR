package main

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"
	"github.com/jimorc/jsdr/internal/ui"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var log *logger.Logger
var mainWin *fyne.Window

func main() {
	logLevel, LogFile := parseCommandLine()

	log = initLogfile(logLevel, LogFile)
	defer log.Close()

	log.Logf(logger.Info, "jsdr started at %v\n", time.Now().UTC())

	a := app.NewWithID("com.github.jimorc.jsdr")
	mainWin := ui.MakeMainWindow(&a, log)
	log.Log(logger.Debug, "Displaying main window\n")
	mainWin.ShowAndRun()
	log.Log(logger.Debug, "Terminated main window\n")
	if ui.SoapyDev.Device != nil {
		sdr.Unmake(ui.SoapyDev, log)
	}

	log.Logf(logger.Info, "jsdr terminated at %v\n", time.Now().UTC())
}

func initLogfile(level logger.LoggingLevel, fileName string) *logger.Logger {
	log, err := logger.NewFileLogger(fileName)
	if err != nil {
		fmt.Printf("Error trying to open log file '%s': %s\n", fileName, err.Error())
		os.Exit(1)
	}
	log.SetMaxLevel(level)
	return log
}

func parseCommandLine() (logger.LoggingLevel, string) {
	pflag.Bool("error", false, "Log fatal and error messages")
	pflag.Bool("info", false, "Log fatal, error, and info messages")
	pflag.Bool("debug", false, "Log fatal, error, info, and debug messages")
	pflag.String("out", os.Getenv("HOME")+"/jsdr.log", "Log filename. If 'stdout', messages are logged to 'stdout'.")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	debug := viper.GetBool("debug")
	info := viper.GetBool("info")
	error := viper.GetBool("error")
	logFile := viper.GetString("out")
	logLevel := logger.Info
	if error {
		logLevel = logger.Error
	}
	if info {
		logLevel = logger.Info
	}
	if debug {
		logLevel = logger.Debug
	}
	return logLevel, logFile
}
