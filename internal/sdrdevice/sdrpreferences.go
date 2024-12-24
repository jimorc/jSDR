package sdrdevice

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"
)

// SdrPreferences holds properties that should be persisted between executions of jsdr.
// For example, this struct holds the last SDR, its last selected sample rate, and its
// last selected antenna.
type SdrPreferences struct {
	Device     string
	SampleRate string
	Antenna    string
}

// variables that are needed across multiple functions and methods.
var jsdrLog *logger.Logger
var parentWindow *fyne.Window
var settingsDialog *dialog.ConfirmDialog
var sdrsSelect *widget.Select
var sampleRatesSelect *widget.Select
var antennaSelect *widget.Select

// ClearPreferences clears all values in the SdrDevice struct. This should only be called if the
// previously stored SDR is no longer connected to the computer.
//
// Params:
//
//	log is the logger to write messages to.
func (s *SdrPreferences) ClearPreferences(log *logger.Logger) {
	s.Device = ""
	s.SampleRate = ""
	s.Antenna = ""
	s.SavePreferences(log)
	log.Log(logger.Debug, "SdrDevice device settings have been cleared\n")
}

// NewFromPreferences creates a new SdrPreferences struct populated with data from
// the program's preferences.
//
// Returns pointer to the new SdrPreferences object.
func NewFromPreferences(log *logger.Logger) *SdrPreferences {
	jsdrLog = log
	s := &SdrPreferences{}
	s.Device = fyne.CurrentApp().Preferences().String("device")
	log.Logf(logger.Debug, "Value: %s loaded from preference: %s\n", s.Device, "device")
	s.SampleRate = fyne.CurrentApp().Preferences().String("samplerate")
	log.Logf(logger.Debug, "Value: %s loaded from preference: %s\n", s.SampleRate, "samplerate")
	s.Antenna = fyne.CurrentApp().Preferences().String("antenna")
	log.Logf(logger.Debug, "Value: %s loaded from preference: %s\n", s.Antenna, "antenna")
	return s
}

// SavePreferences saves the values in the SdrPreferences object to the program's preferences
// file.
//
// SdrC
// Params:
//
//	log - the logger to write log messages to.
func (s *SdrPreferences) SavePreferences(log *logger.Logger) {
	savePreference(s.Device, "device", log)
	savePreference(s.SampleRate, "samplerate", log)
	savePreference(s.Antenna, "antenna", log)
}

// MakeSettingsAction creates a ToolbarAction widget linked to showing the settings dialog.
func (s *SdrPreferences) MakeSettingsAction() *widget.ToolbarAction {
	return widget.NewToolbarAction(theme.SettingsIcon(), s.showSettingsDialog)
}

// CreateSettingsDialog creates the settings dialog.
//
// Params:
//
//	parent is the parent window for this dialog. This should probably be the main window.
//	log is thelogger to write log messages to.
func (s *SdrPreferences) CreateSettingsDialog(parent *fyne.Window, log *logger.Logger) {
	parentWindow = parent
	log.Log(logger.Debug, "In settingsCallback\n")
	sdrsLabel := widget.NewLabel("SDR Device:")
	sdrsSelect = widget.NewSelect([]string{}, s.SdrChanged)

	sampleRateLabel := widget.NewLabel("Sample Rate:")
	sampleRateLabel.Alignment = fyne.TextAlignTrailing
	sampleRatesSelect = widget.NewSelect([]string{}, s.SampleRateChanged)

	antennaLabel := widget.NewLabel("Antenna:")
	antennaLabel.Alignment = fyne.TextAlignTrailing
	antennaSelect = widget.NewSelect([]string{}, s.AntennaChanged)
	grid := container.NewGridWithColumns(2, sdrsLabel, sdrsSelect, sampleRateLabel, sampleRatesSelect,
		antennaLabel, antennaSelect)
	settingsDialog = dialog.NewCustomConfirm("SDR Settings", "Accept", "Close", grid, acceptCancelCallback, *parentWindow)
}

// showSettingsDialog is the callback for the settings ToolbarAction widget created in MakeSettingsAction.
func (s *SdrPreferences) showSettingsDialog() {
	// retrieve the list of attached SDRs.
	sdrs := sdr.EnumerateSdrsWithoutAudio(sdr.SoapyDev, jsdrLog)
	jsdrLog.Logf(logger.Debug, "Number of sdr devices returned from EnumerateSdrsWithoutAudio: %d\n",
		sdrs.NumberOfSdrs())
	if sdrs.NumberOfSdrs() == 0 {
		noDevices := dialog.NewInformation("No Attached SDRs",
			"No SDRs were found.\nAttach an SDR, then try again.",
			*parentWindow)
		noDevices.Show()
		return
	} else {
		// there is at least one SDR, so retrieve their labels and set as Options in sdrsSelect.
		sdrLabels := sdrs.SdrLabels()
		sdrsSelect.Options = sdrLabels
		if len(sdrLabels) == 1 {
			sdrsSelect.SetSelectedIndex(0)
		} else {
			sdrsSelect.SetSelected(s.Device)
		}
		settingsDialog.Show()
	}
}

// savePreference saves the specified value to the name preference.
//
// Params:
//
//	pref - the string to be saved.
//	prefName - the preference name to save to.
//	log - the logger to write log messages to.
func savePreference(pref string, prefName string, log *logger.Logger) {
	log.Logf(logger.Debug, "Value: %s saved to preference: %s\n", pref, prefName)
	fyne.CurrentApp().Preferences().SetString(prefName, pref)
}

// SdrChanged is the callback executed when an SDR is selected in the settings dialog.
func (s SdrPreferences) SdrChanged(selectedSdr string) {
	jsdrLog.Logf(logger.Debug, "SDR selected: %s\n", selectedSdr)

}

// SampleRateChanged is the callback executed when one of the sample rates is selected
// in the settings dialog.
func (sP SdrPreferences) SampleRateChanged(sampleRate string) {
	jsdrLog.Logf(logger.Debug, "Sample Rate selected: %s\n", sampleRate)
}

// AntennaChanged is the callback executed when one of the antennas is selected in
// the settings dialog.
func (sP SdrPreferences) AntennaChanged(antenna string) {
	jsdrLog.Logf(logger.Debug, "Antenna selected: %s\n", antenna)
}

// acceptCancelCallback is the callback that is called when either of the Close or Accept
// buttons in the settings dialog is clicked.
// When called, the fyne framework closes the settings dialog.
//
// Params:
//
//	accept indicates which of the Close or Accept buttons was clicked:
//		true indicates that the Accept button was clicked.
//		false indicates that the Close button was clicked.
func acceptCancelCallback(accept bool) {
	jsdrLog.Logf(logger.Debug, "In settingsDialogCallback: %v\n", accept)
}
