package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"
)

// sdrs is a map of devices info indexed by device's label
var sdrs map[string]map[string]string

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
		sdrsSelect := widget.NewSelect(sdrLabels, sdrChanged)
		grid := container.NewGridWithColumns(2, widget.NewLabel("SDR Device:"), sdrsSelect)
		settings := dialog.NewCustomConfirm("SDR Settings", "Accept", "Close", grid, settingsDialogCallback, mainWin)
		settings.Show()
		if len(sdrLabels) == 1 {
			sdrsSelect.SetSelectedIndex(0)
		}
	}
}

func settingsDialogCallback(accept bool) {
	jsdrLogger.Logf(logger.Debug, "In settingsDialogCallback: %v\n", accept)
}

func sdrChanged(value string) {
	jsdrLogger.Logf(logger.Debug, "SDR selected: %s\n", value)
}
