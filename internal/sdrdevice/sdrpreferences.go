package sdrdevice

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/jimorc/jsdr/internal/logger"
)

// SdrPreferences holds properties that should be persisted between executions of jsdr.
// For example, this struct holds the last SDR, its last selected sample rate, and its
// last selected antenna.
type SdrPreferences struct {
	Device     binding.String
	SampleRate binding.String
	Antenna    binding.String
}

type SdrChanges interface {
	SdrChanged()
	SampleRateChanged()
	AntennaChanged()
}

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
func NewFromPreferences(sP SdrChanges) *SdrPreferences {
	s := &SdrPreferences{}
	s.Device = bindToString("device", sP.SdrChanged)
	s.SampleRate = bindToString("samplerate", sP.SampleRateChanged)
	s.Antenna = bindToString("antenna", sP.AntennaChanged)
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

// bindToString creates a binding to the specified string, loads the specified value
// from the program's preferences, and adds a listener.
//
// Params:
//
//	prefName - the name of the value to retrieve from the program's preferences.
//
//	listener - the callback function to be called whenever the value in the bound
//
// string is modified.
//
// Returns the bound string object.
func bindToString(prefName string, listener func()) binding.String {
	s := binding.NewString()
	s.Set(fyne.CurrentApp().Preferences().String(prefName))
	callback := binding.NewDataListener(listener)
	s.AddListener(callback)
	return s
}

func (sP SdrPreferences) SdrChanged() {

}

func (sP SdrPreferences) SampleRateChanged() {

}

func (sP SdrPreferences) AntennaChanged() {

}
