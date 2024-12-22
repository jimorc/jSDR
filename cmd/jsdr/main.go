package main

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"
	"github.com/jimorc/jsdr/internal/sdrdevice"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var log *logger.Logger
var mainWin fyne.Window
var sdrs sdr.Sdrs
var sdrDevice sdrdevice.SdrDevice
var sampleRatesSelect *widget.Select
var antennaSelect *widget.Select

func main() {
	logLevel, LogFile := parseCommandLine()

	log = initLogfile(logLevel, LogFile)
	defer log.Close()

	log.Logf(logger.Info, "jsdr started at %v\n", time.Now().UTC())

	a := app.NewWithID("com.github.jimorc.jsdr")
	sdrDevice.LoadFromApp(log)
	mainWin = makeMainWindow(&a, log)

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

func makeMainWindow(jsdrApp *fyne.App, log *logger.Logger) fyne.Window {
	mainWin := (*jsdrApp).NewWindow("jsdr")

	log.Log(logger.Debug, "Creating main window content\n")
	settingsAction := makeSettingsAction()
	toolbar := widget.NewToolbar(settingsAction)
	mainWin.SetContent(container.NewBorder(toolbar, nil, nil, nil))
	mainWin.Resize(fyne.NewSize(800, 400))
	log.Log(logger.Debug, "Main window content created\n")
	return mainWin
}

func makeSettingsAction() *widget.ToolbarAction {
	return widget.NewToolbarAction(theme.SettingsIcon(), settingsCallback)
}

func settingsCallback() {
	log.Log(logger.Debug, "In settingsCallback\n")
	sdrs = sdr.EnumerateSdrsWithoutAudio(sdr.SoapyDev, log)
	log.Logf(logger.Debug, "Number of sdr devices returned from EnumerateSdrsWithoutAudio: %d\n",
		sdrs.NumberOfSdrs())
	if sdrs.NumberOfSdrs() == 0 {
		noDevices := dialog.NewInformation("No Attached SDRs",
			"No SDRs were found.\nAttach an SDR, then try again.",
			mainWin)
		noDevices.Show()
	} else {
		var sdrLabels []string
		for k := range sdrs.DevicesMap {
			sdrLabels = append(sdrLabels, k)
		}
		sdrsLabel := widget.NewLabel("SDR Device:")
		sdrsSelect := widget.NewSelect(sdrLabels, sdrChanged)
		sampleRateLabel := widget.NewLabel("Sample Rate:")
		sampleRateLabel.Alignment = fyne.TextAlignTrailing
		sampleRatesSelect = widget.NewSelect([]string{}, sampleRateChanged)
		antennaLabel := widget.NewLabel("Antenna:")
		antennaLabel.Alignment = fyne.TextAlignTrailing
		antennaSelect = widget.NewSelect([]string{}, antennaChanged)
		grid := container.NewGridWithColumns(2, sdrsLabel, sdrsSelect, sampleRateLabel, sampleRatesSelect,
			antennaLabel, antennaSelect)
		settings := dialog.NewCustomConfirm("SDR Settings", "Accept", "Close", grid, settingsDialogCallback, mainWin)
		settings.Show()
		if len(sdrLabels) == 1 {
			sdrsSelect.SetSelectedIndex(0)
		}
	}
}

func antennaChanged(antenna string) {
	log.Logf(logger.Debug, "Antenna selected: %s\n", antenna)
}

func settingsDialogCallback(accept bool) {
	log.Logf(logger.Debug, "In settingsDialogCallback: %v\n", accept)
}

func sdrChanged(value string) {
	log.Logf(logger.Debug, "SDR selected: %s\n", value)
	devProps := sdrs.DevicesMap[value]
	if sdr.SoapyDev.Device != nil {
		sdr.Unmake(sdr.SoapyDev, log)
	}
	err := sdr.Make(sdr.SoapyDev, devProps, log)
	if err != nil {
		errDialog := dialog.NewError(err, mainWin)
		errDialog.Show()
	} else {

		sampleRatesSelect.Options = sdr.GetSampleRates(sdr.SoapyDev, log)
		sampleRatesSelect.Selected = sdr.GetSampleRate(sdr.SoapyDev, log)
		sampleRatesSelect.Refresh()
		antennaSelect.Options = sdr.GetAntennaNames(sdr.SoapyDev, log)
		if len(antennaSelect.Options) == 1 {
			antennaSelect.SetSelectedIndex(0)
		} else {
			antennaSelect.SetSelected(sdr.GetCurrentAntenna(sdr.SoapyDev, log))
		}
		antennaSelect.Refresh()
	}
}

func sampleRateChanged(rate string) {
	log.Logf(logger.Debug, "Sample rate selected: %s\n", rate)
}
