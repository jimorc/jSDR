package sdrdevice

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
	Device     binding.String
	SampleRate binding.String
	Antenna    binding.String
}

// variables that are needed across multiple functions and methods.
var jsdrLog *logger.Logger
var parentWindow *fyne.Window
var settingsDialog *dialog.ConfirmDialog
var sdrsSelect *widget.Select

// ClearPreferences clears all values in the SdrDevice struct. This should only be called if the
// previously stored SDR is no longer connected to the computer.
//
// Params:
//
//	log is the logger to write messages to.
func (s *SdrPreferences) ClearPreferences(log *logger.Logger) error {
	s.Device.Set("")
	s.SampleRate.Set("")
	s.Antenna.Set("")
	err := s.SavePreferences(log)
	if err != nil {
		log.Logf(logger.Debug, "Error encountered while clearing SDR Preferences: %s.\n",
			err.Error())
	} else {
		log.Log(logger.Debug, "SdrDevice device settings have been cleared\n")
	}
	return err
}

// NewFromPreferences creates a new SdrPreferences struct populated with data from
// the program's preferences.
//
// Returns pointer to the new SdrPreferences object.
func NewFromPreferences(log *logger.Logger) *SdrPreferences {
	jsdrLog = log
	s := &SdrPreferences{}
	s.Device = binding.BindPreferenceString("device", fyne.CurrentApp().Preferences())
	s.SampleRate = binding.BindPreferenceString("samplerate", fyne.CurrentApp().Preferences())
	s.Antenna = binding.BindPreferenceString("antenna", fyne.CurrentApp().Preferences())
	return s
}

// SavePreferences saves the values in the SdrPreferences object to the program's preferences
// file.
//
// Params:
//
//	log - the logger to write log messages to.
func (s *SdrPreferences) SavePreferences(log *logger.Logger) error {
	err := savePreference(s.Device, "device", log)
	if err != nil {
		return err
	}
	err = savePreference(s.SampleRate, "samplerate", log)
	if err != nil {
		return err
	}
	err = savePreference(s.Antenna, "antenna", log)
	if err != nil {
		return err
	}
	return nil
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
	sdrsSelect.Bind(s.Device)

	sampleRateLabel := widget.NewLabel("Sample Rate:")
	sampleRateLabel.Alignment = fyne.TextAlignTrailing
	sampleRatesSelect := widget.NewSelect([]string{}, s.SampleRateChanged)
	sampleRatesSelect.Bind(s.SampleRate)

	antennaLabel := widget.NewLabel("Antenna:")
	antennaLabel.Alignment = fyne.TextAlignTrailing
	antennaSelect := widget.NewSelect([]string{}, s.AntennaChanged)
	antennaSelect.Bind(s.Antenna)
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
		// there is at least one SDR, so
		var sdrLabels []string
		for k := range sdrs.DevicesMap {
			sdrLabels = append(sdrLabels, k)
		}

		sdrsSelect.Options = sdrLabels
		if len(sdrLabels) == 1 {
			s.Device.Set(sdrLabels[0])
		}
		settingsDialog.Show()
	}
}

// savePreference saves the specified value to the name preference.
//
// Params:
//
//	pref - the bound string to be saved.
//	prefName - the preference name to save to.
//	log - the logger to write log messages to.
func savePreference(pref binding.String, prefName string, log *logger.Logger) error {
	value, err := pref.Get()
	if err != nil {
		log.Logf(logger.Error, "Unable to retrieve %s value.\n", prefName)
		return err
	} else {
		fyne.CurrentApp().Preferences().SetString(prefName, value)
		return nil
	}

}

// SdrChanged is the callback executed when an SDR is selected in the settings dialog.
func (sP SdrPreferences) SdrChanged(selectedSdr string) {
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
