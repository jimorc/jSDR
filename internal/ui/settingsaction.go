package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"
)

// sdrs is a map of devices info indexed by device's label
var sdrs sdr.Sdrs
var sampleRatesSelect *widget.Select
var antennaSelect *widget.Select

func makeSettingsAction() *widget.ToolbarAction {
	JsdrLogger.Log(logger.Debug, "Entered ui.makeSettingsAction\n")
	action := widget.NewToolbarAction(theme.SettingsIcon(), settingsCallback)
	JsdrLogger.Log(logger.Debug, "Returning the settings toolbar action from makeSettingsAction\n")
	return action
}

func settingsCallback() {
	JsdrLogger.Log(logger.Debug, "In settingsCallback\n")
	sdrs = sdr.EnumerateSdrsWithoutAudio(sdr.SoapyDev, JsdrLogger)
	JsdrLogger.Logf(logger.Debug, "Number of sdr devices returned from EnumerateSdrsWithoutAudio: %d\n",
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
	JsdrLogger.Logf(logger.Debug, "Antenna selected: %s\n", antenna)
}

func settingsDialogCallback(accept bool) {
	JsdrLogger.Logf(logger.Debug, "In settingsDialogCallback: %v\n", accept)
}

func sdrChanged(value string) {
	JsdrLogger.Logf(logger.Debug, "SDR selected: %s\n", value)
	devProps := sdrs.DevicesMap[value]
	if sdr.SoapyDev.Device != nil {
		sdr.Unmake(sdr.SoapyDev, JsdrLogger)
	}
	err := sdr.Make(sdr.SoapyDev, devProps, JsdrLogger)
	if err != nil {
		errDialog := dialog.NewError(err, mainWin)
		errDialog.Show()
	} else {

		sampleRatesSelect.Options = sdr.GetSampleRates(sdr.SoapyDev, JsdrLogger)
		sampleRatesSelect.Selected = sdr.GetSampleRate(sdr.SoapyDev, JsdrLogger)
		sampleRatesSelect.Refresh()
		antennaSelect.Options = sdr.GetAntennaNames(sdr.SoapyDev, JsdrLogger)
		if len(antennaSelect.Options) == 1 {
			antennaSelect.SetSelectedIndex(0)
		} else {
			antennaSelect.SetSelected(sdr.GetCurrentAntenna(sdr.SoapyDev, JsdrLogger))
		}
		antennaSelect.Refresh()
	}
}

func sampleRateChanged(rate string) {
	JsdrLogger.Logf(logger.Debug, "Sample rate selected: %s\n", rate)
}
