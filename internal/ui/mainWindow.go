package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jimorc/jsdr/internal/logger"
)

func MakeMainWindow(jsdrApp *fyne.App, log *logger.Logger) fyne.Window {
	win := (*jsdrApp).NewWindow("jsdr")

	log.Log(logger.Debug, "Creating main window content\n")
	settingsAction := widget.NewToolbarAction(theme.SettingsIcon(), func() {
		log.Log(logger.Debug, "Inside settingsAction callback\n")
	})
	toolbar := widget.NewToolbar(settingsAction)
	win.SetContent(container.NewBorder(toolbar, nil, nil, nil))
	win.Resize(fyne.NewSize(400, 200))
	log.Log(logger.Debug, "Main window content created\n")
	return win
}
