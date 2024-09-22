package main

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jimorc/jsdr/internal/logger"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	logLevel, logFile := parseCommandLine()

	log := initLogfile(logLevel, logFile)
	defer log.Close()

	log.Logf(logger.Info, "jsdr started at %v\n", time.Now().UTC())

	a := app.NewWithID("com.github.jimorc.jsdr")
	win := a.NewWindow("jsdr")
	win.SetContent(makeUI(&win, log))
	win.Resize(fyne.NewSize(400, 200))
	log.Log(logger.Debug, "Displaying main window\n")
	win.ShowAndRun()
	log.Log(logger.Debug, "Terminated main window\n")

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

func makeUI(_ *fyne.Window, log *logger.Logger) fyne.CanvasObject {
	log.Log(logger.Debug, "Creating main window content\n")
	settingsAction := widget.NewToolbarAction(theme.SettingsIcon(), func() {
		log.Log(logger.Debug, "Inside settingsAction callback\n")
	})
	toolbar := widget.NewToolbar(settingsAction)
	winContainer := container.NewBorder(toolbar, nil, nil, nil)
	log.Log(logger.Debug, "Main window content created\n")
	return winContainer
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
