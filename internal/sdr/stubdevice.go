package sdr

// StubDevice provides a stub interface for testing soapySDR interface.
// The fields in the struct allow loading of data for test purposesto allow simple method returns.
type StubDevice struct {
	Devices []map[string]string
}

// Enumerate returns a slice of map[string]string values representing the available devices. These
// values must be preloaded into the Devices property of the StubDevice struct before Enumerate is called.
func (dev StubDevice) Enumerate(args map[string]string) []map[string]string {
	return dev.Devices
}
