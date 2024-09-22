package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/jimorc/jsdr/internal/logger"
)

func makeSettingsAction() *widget.ToolbarAction {
	uiLogger.Log(logger.Debug, "Entered ui.makeSettingsAction\n")
	action := widget.NewToolbarAction(theme.SettingsIcon(), settingsCallback)
	uiLogger.Log(logger.Debug, "Returning the settings toolbar action from makeSettingsAction\n")
	return action
}

func settingsCallback() {
	uiLogger.Log(logger.Debug, "In settingsCallback\n")
	label := widget.NewLabel("Settings")
	content := container.NewVBox(label)
	settings := dialog.NewCustomConfirm("SDR Settings", "Accept", "Close", content, settingsDialogCallback, mainWin)
	settings.Show()
}

func settingsDialogCallback(accept bool) {
	uiLogger.Logf(logger.Debug, "In settingsDialogCallback: %v\n", accept)
}
