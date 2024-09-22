package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jimorc/jsdr/internal/logger"
)

var uiLogger *logger.Logger
var mainWin fyne.Window

func MakeMainWindow(jsdrApp *fyne.App, log *logger.Logger) fyne.Window {
	uiLogger = log
	mainWin = (*jsdrApp).NewWindow("jsdr")

	uiLogger.Log(logger.Debug, "Creating main window content\n")
	settingsAction := makeSettingsAction()
	toolbar := widget.NewToolbar(settingsAction)
	mainWin.SetContent(container.NewBorder(toolbar, nil, nil, nil))
	mainWin.Resize(fyne.NewSize(400, 200))
	uiLogger.Log(logger.Debug, "Main window content created\n")
	return mainWin
}
