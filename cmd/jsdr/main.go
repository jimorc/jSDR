package main

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"
	"github.com/jimorc/jsdr/internal/sdrdevice"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// variables that are needed across multiple functions.
var log *logger.Logger
var mainWin fyne.Window
var sdrPrefs sdrdevice.SdrPreferences

func main() {
	logLevel, LogFile := parseCommandLine()

	log = initLogfile(logLevel, LogFile)
	defer log.Close()

	log.Logf(logger.Info, "jsdr started at %v\n", time.Now().UTC())
	a := app.NewWithID("com.github.jimorc.jsdr")
	sdrPrefs = *sdrdevice.NewFromPreferences(log)
	defer sdrPrefs.SavePreferences(log)
	mainWin = makeMainWindow(&a, &sdrPrefs, log)
	sdrPrefs.CreateSettingsDialog(&mainWin, log)

	log.Log(logger.Debug, "Displaying main window\n")
	mainWin.ShowAndRun()
	log.Log(logger.Debug, "Terminated main window\n")
	if sdr.SoapyDev.Device != nil {
		sdr.Unmake(sdr.SoapyDev, log)
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

func makeMainWindow(jsdrApp *fyne.App, prefs *sdrdevice.SdrPreferences, log *logger.Logger) fyne.Window {
	mainWin := (*jsdrApp).NewWindow("jsdr")

	log.Log(logger.Debug, "Creating main window content\n")
	settingsAction := prefs.MakeSettingsAction()
	toolbar := widget.NewToolbar(settingsAction)
	mainWin.SetContent(container.NewBorder(toolbar, nil, nil, nil))
	mainWin.Resize(fyne.NewSize(800, 400))
	log.Log(logger.Debug, "Main window content created\n")
	return mainWin
}
