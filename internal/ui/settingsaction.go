package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"
)

func makeSettingsAction() *widget.ToolbarAction {
	jsdrLogger.Log(logger.Debug, "Entered ui.makeSettingsAction\n")
	action := widget.NewToolbarAction(theme.SettingsIcon(), settingsCallback)
	jsdrLogger.Log(logger.Debug, "Returning the settings toolbar action from makeSettingsAction\n")
	return action
}

func settingsCallback() {
	jsdrLogger.Log(logger.Debug, "In settingsCallback\n")
	sdrs := widget.NewSelect(sdr.EnumerateWithoutAudio(jsdrLogger), sdrChanged)
	grid := container.NewGridWithColumns(2, widget.NewLabel("SDR Device:"), sdrs)
	settings := dialog.NewCustomConfirm("SDR Settings", "Accept", "Close", grid, settingsDialogCallback, mainWin)
	settings.Show()
}

func settingsDialogCallback(accept bool) {
	jsdrLogger.Logf(logger.Debug, "In settingsDialogCallback: %v\n", accept)
}

func sdrChanged(value string) {
	jsdrLogger.Logf(logger.Debug, "SDR selected: %s\n", value)
}
