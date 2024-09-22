package ui

import (
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/jimorc/jsdr/internal/logger"
)

func makeSettingsAction() *widget.ToolbarAction {
	uiLogger.Log(logger.Debug, "Entered ui.makeSettingsAction\n")
	action := widget.NewToolbarAction(theme.SettingsIcon(), func() {
		uiLogger.Log(logger.Debug, "Inside settingsAction callback\n")
	})
	uiLogger.Log(logger.Debug, "Returning the settings toolbar action from makeSettingsAction\n")
	return action
}
