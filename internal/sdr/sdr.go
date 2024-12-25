// Package sdr provides interfaces and functions that allow multiple device types.
//
// The initial device types are:
//
//	SoapyDevice for SoapySDR devices.
//
//	StubDevice for testing of the various sdr functions.
//
// Many of the function and method names are changed from those provided in go-soapy-sdr.go.
// I find many of the function and method names to be confusing in go-soapy-sdr.go For example:
// device.SetAntennas sets a single antenna on a device, not multiple antennas.
//
// These name changes are an attempt to clarify what the functions and methods do,
package sdr

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

// MakeDevice interface specifies the methods for creating and destroying an SDR device.
type MakeDevice interface {
	Make(args map[string]string) error
	Unmake() error
	GetHardwareKey() string
}

// KeyValues interface specifies the methods for retrieving SDR information.
type KeyValues interface {
	GetHardwareKey() string
}

// Sdr represents the SDR device.
type Sdr struct {
	Device           *device.SDRDevice
	DeviceName       string
	DeviceNames      []string
	DeviceProperties map[string]string
	SampleRates      []string
	SampleRate       float64
	Antennas         []string
	Antenna          string
}

// variables that are needed across multiple functions and methods.
var jsdrLog *logger.Logger
var SoapyDev = &SoapyDevice{}
var parentWindow *fyne.Window
var settingsDialog *dialog.ConfirmDialog
var sdrsSelect *widget.Select
var sampleRatesSelect *widget.Select
var antennaSelect *widget.Select

// Make makes a new device given construction args.
//
// Construction args should be as explicit as possible (i.e. include all values retrieved by
// EnumerateSdrsWithoutAudio). args should contain a label value.
func Make(sdrD MakeDevice, args map[string]string, log *logger.Logger) error {
	log.Logf(logger.Debug, "Making device with label: %s\n", args["label"])
	err := sdrD.Make(args)
	if err != nil {
		log.Logf(logger.Error, "Error encountered trying to make device: %s\n", err.Error())
		return err
	}
	log.Logf(logger.Debug, "Made SDR with hardware key: %s\n", sdrD.GetHardwareKey())
	return nil
}

// Unmake frees up any assets associated with the SDR device.
//
// No sdr calls should be made after Unmake is called.
func Unmake(sdrD MakeDevice, log *logger.Logger) error {
	log.Log(logger.Debug, "Attempting to unmake an SDR device\n")
	err := sdrD.Unmake()
	if err != nil {
		log.Logf(logger.Error, "Error attempting to unmake an SDR device: %s\n", err.Error())
	}
	return err
}

// GetHardwareKey returns the hardware key for the SDR device.
func (sdr *Sdr) GetHardwareKey(sdrD KeyValues) string {
	return sdrD.GetHardwareKey()
}

// ClearPreferences clears all values in the SdrDevice struct. This should only be called if the
// previously stored SDR is no longer connected to the computer.
//
// Params:
//
//	log is the logger to write messages to.
func (sdr *Sdr) ClearPreferences(log *logger.Logger) {
	sdr.DeviceName = ""
	sdr.SampleRate = 0.
	sdr.Antenna = ""
	sdr.SavePreferences(log)
	log.Log(logger.Debug, "Sdr device settings have been cleared\n")
}

// LoadPreferences loads preferences values into an Sdr object from
// the program's preferences.
//
// Returns pointer to the new SdrPreferences object.

func (sdr *Sdr) LoadPreferences(log *logger.Logger) {
	sdr.DeviceName = fyne.CurrentApp().Preferences().String("device")
	log.Logf(logger.Debug, "Value: %s loaded from preference: %s\n", sdr.DeviceName, "device")
	sdr.SampleRate = fyne.CurrentApp().Preferences().Float("samplerate")
	log.Logf(logger.Debug, "Value: %f loaded from preference: %s\n", sdr.SampleRate, "samplerate")
	sdr.Antenna = fyne.CurrentApp().Preferences().String("antenna")
	log.Logf(logger.Debug, "Value: %s loaded from preference: %s\n", sdr.Antenna, "antenna")
}

// SavePreferences saves the values in the Sdr object to the program's preferences
// file.
//
// Params:
//
//	log - the logger to write log messages to.
func (sdr *Sdr) SavePreferences(log *logger.Logger) {
	saveStringPreference(sdr.DeviceName, "device", log)
	saveFloatPreference(sdr.SampleRate, "samplerate", log)
	saveStringPreference(sdr.Antenna, "antenna", log)
}

// saveStringPreference saves the specified value to the name preference.
//
// Params:
//
//	pref - the string to be saved.
//	prefName - the preference name to save to.
//	log - the logger to write log messages to.
func saveStringPreference(pref string, prefName string, log *logger.Logger) {
	log.Logf(logger.Debug, "Value: %s saved to preference: %s\n", pref, prefName)
	fyne.CurrentApp().Preferences().SetString(prefName, pref)
}

// saveFloatPreference saves the specified value to the name preference.
//
// Params:
//
//	pref - the float to be saved.
//	prefName - the preference name to save to.
//	log - the logger to write log messages to.
func saveFloatPreference(pref float64, prefName string, log *logger.Logger) {
	log.Logf(logger.Debug, "Value: %f saved to preference: %s\n", pref, prefName)
	fyne.CurrentApp().Preferences().SetFloat(prefName, pref)
}

// MakeSettingsAction creates a ToolbarAction widget linked to showing the settings dialog.
func (sdr *Sdr) MakeSettingsAction() *widget.ToolbarAction {
	return widget.NewToolbarAction(theme.SettingsIcon(), sdr.showSettingsDialog)
}

// CreateSettingsDialog creates the settings dialog.
//
// Params:
//
//	parent is the parent window for this dialog. This should probably be the main window.
//	log is thelogger to write log messages to.
func (sdr *Sdr) CreateSettingsDialog(parent *fyne.Window, log *logger.Logger) {
	jsdrLog = log
	parentWindow = parent
	log.Log(logger.Debug, "In settingsCallback\n")
	sdrsLabel := widget.NewLabel("SDR Device:")
	sdrsLabel.Alignment = fyne.TextAlignTrailing
	sdrsSelect = widget.NewSelect([]string{}, sdr.SdrChanged)

	sampleRateLabel := widget.NewLabel("Sample Rate:")
	sampleRateLabel.Alignment = fyne.TextAlignTrailing
	sampleRatesSelect = widget.NewSelect([]string{}, sdr.SampleRateChanged)

	antennaLabel := widget.NewLabel("Antenna:")
	antennaLabel.Alignment = fyne.TextAlignTrailing
	antennaSelect = widget.NewSelect([]string{}, sdr.AntennaChanged)
	grid := container.New(layout.NewFormLayout(), sdrsLabel, sdrsSelect, sampleRateLabel, sampleRatesSelect,
		antennaLabel, antennaSelect)
	settingsDialog = dialog.NewCustomConfirm("SDR Settings", "Accept", "Close", grid, acceptCancelCallback, *parentWindow)
	settingsDialog.Resize(fyne.NewSize(450, settingsDialog.MinSize().Height))
}

// showSettingsDialog is the callback for the settings ToolbarAction widget created in MakeSettingsAction.
func (sdr *Sdr) showSettingsDialog() {
	// retrieve the list of attached SDRs.
	sdrs := EnumerateSdrsWithoutAudio(SoapyDev, jsdrLog)
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
		sdr.DeviceNames = sdrLabels
		sdrsSelect.Options = sdrLabels
		if len(sdrLabels) == 1 {
			sdrsSelect.SetSelectedIndex(0)
		} else {
			sdrsSelect.SetSelected(sdr.DeviceName)
		}
		settingsDialog.Show()
	}
}

// SdrChanged is the callback executed when an SDR is selected in the settings dialog.
func (sdr *Sdr) SdrChanged(selectedSdr string) {
	jsdrLog.Logf(logger.Debug, "SDR selected: %s\n", selectedSdr)
	sdrs := EnumerateSdrsWithoutAudio(SoapyDev, jsdrLog)
	s := sdrs.Sdr(selectedSdr)
	if s == nil {
		jsdrLog.Logf(logger.Error, "The selected SDR: %s is no longer attached.\n", selectedSdr)
		settingsDialog.Hide()
		noDevice := dialog.NewInformation("SDR No Longer Attached",
			"The selected SDR is no longer attached to the computer.\n"+
				"Either reattach the SDR and try again, or select another SDR.",
			*parentWindow)
		noDevice.Show()
		return
	}
	sdr.DeviceProperties = s
}

// SampleRateChanged is the callback executed when one of the sample rates is selected
// in the settings dialog.
func (sdr *Sdr) SampleRateChanged(sampleRate string) {
	jsdrLog.Logf(logger.Debug, "Sample Rate selected: %s\n", sampleRate)
}

// AntennaChanged is the callback executed when one of the antennas is selected in
// the settings dialog.
func (sdr *Sdr) AntennaChanged(antenna string) {
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
