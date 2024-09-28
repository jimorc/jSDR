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

var sdrsSelect *widget.Select

// sdrs is a map of devices info indexed by device's label
var sdrs map[string]map[string]string
var selSdr *sdr.Sdr
var sampleRatesSelect *widget.Select
var antennaSelect *widget.Select

func makeSettingsAction() *widget.ToolbarAction {
	jsdrLogger.Log(logger.Debug, "Entered ui.makeSettingsAction\n")
	action := widget.NewToolbarAction(theme.SettingsIcon(), settingsCallback)
	jsdrLogger.Log(logger.Debug, "Returning the settings toolbar action from makeSettingsAction\n")
	return action
}

func settingsCallback() {
	jsdrLogger.Log(logger.Debug, "In settingsCallback\n")
	sdrs = sdr.EnumerateWithoutAudio(jsdrLogger)
	jsdrLogger.Logf(logger.Debug, "Number of sdr devices returned from EnumerateWithoutAudio: %d\n", len(sdrs))
	if len(sdrs) == 0 {
		noDevices := dialog.NewInformation("No Attached SDRs",
			"No SDRs were found.\nAttach an SDR, then try again.",
			mainWin)
		noDevices.Show()
	} else {
		var sdrLabels []string
		for k := range sdrs {
			sdrLabels = append(sdrLabels, k)
		}
		sdrsLabel := widget.NewLabel("SDR Device:")
		sdrsLabel.Alignment = fyne.TextAlignTrailing
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
	jsdrLogger.Logf(logger.Debug, "Antenna selected: %s\n", antenna)
}

func settingsDialogCallback(accept bool) {
	jsdrLogger.Logf(logger.Debug, "In settingsDialogCallback: %v\n", accept)
}

func sdrChanged(value string) {
	jsdrLogger.Logf(logger.Debug, "SDR selected: %s\n", value)
	devProps := sdrs[value]
	dev, err := sdr.Make(devProps, jsdrLogger)
	if err != nil {
		errDialog := dialog.NewError(err, mainWin)
		errDialog.Show()
	} else {
		sampleRatesSelect.Options = dev.GetSampleRates(jsdrLogger)
		sampleRatesSelect.Selected = dev.GetSampleRate(jsdrLogger)
		antennaSelect.Options = dev.GetAntennas(jsdrLogger)
		if len(antennaSelect.Options) == 1 {
			antennaSelect.SetSelectedIndex(0)
		} else {
			antennaSelect.SetSelected(dev.GetCurrentAntenna(jsdrLogger))
		}
	}
}

func sampleRateChanged(rate string) {
	jsdrLogger.Logf(logger.Debug, "Sample rate selected: %s\n", rate)
}
