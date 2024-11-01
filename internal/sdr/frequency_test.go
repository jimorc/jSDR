package sdr_test

import (
	"testing"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/jimorc/jsdr/internal/sdr"
	"github.com/stretchr/testify/assert"
)

func TestGetFrequencyRanges_OneRange(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	ranges, err := sdr.GetFrequencyRanges(&stub, testLogger)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ranges))
	assert.Equal(t, 0.0, ranges[0].Minimum)
	assert.Equal(t, 6e+09, ranges[0].Maximum)
	assert.Equal(t, 0.0, ranges[0].Step)
}

func TestGetFrequencyRanges_TwoRanges(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "3"}}
	ranges, err := sdr.GetFrequencyRanges(&stub, testLogger)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(ranges))
	assert.Equal(t, 0.0, ranges[0].Minimum)
	assert.Equal(t, 6e+09, ranges[0].Maximum)
	assert.Equal(t, 0.0, ranges[0].Step)
	assert.Equal(t, 6.1e+09, ranges[1].Minimum)
	assert.Equal(t, 1e+10, ranges[1].Maximum)
	assert.Equal(t, 0.0, ranges[1].Step)
}

func TestGetFrequencyRanges_NoRanges(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	ranges, err := sdr.GetFrequencyRanges(&stub, testLogger)
	assert.NotNil(t, err)
	assert.Equal(t, "The attached SDR seems defective; there are no specified frequency ranges.", err.Error())
	assert.Equal(t, 0, len(ranges))
}

func TestGetTunableElementNames(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	tElts := sdr.GetTunableElementNames(&stub, testLogger)
	assert.Equal(t, 1, len(tElts))
	assert.Equal(t, "RF", tElts[0])
}

func TestGetTunableElementsFrequencyRanges(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	fRanges, err := sdr.GetTunableElementFrequencyRanges(&stub, testLogger, "RF")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(fRanges))
	assert.Equal(t, 0.0, fRanges[0].Minimum)
	assert.Equal(t, 6e+09, fRanges[0].Maximum)
	assert.Equal(t, 0.0, fRanges[0].Step)
	assert.Equal(t, 6.1e+09, fRanges[1].Minimum)
	assert.Equal(t, 1e+10, fRanges[1].Maximum)
	assert.Equal(t, 0.0, fRanges[1].Step)
}

func TestGetTunableElementsFrequencyRanges_BadElement(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	fRanges, err := sdr.GetTunableElementFrequencyRanges(&stub, testLogger, "IF")
	assert.NotNil(t, err)
	assert.Equal(t, "Invalid tunable element name: IF", err.Error())
	assert.Equal(t, 0, len(fRanges))
}

func TestGetTunableElementFrequency(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	freq, err := sdr.GetTunableElementFrequency(&stub, testLogger, "RF")
	assert.Nil(t, err)
	assert.Equal(t, 1e+08, freq)
}

func TestGetTunableElementFrequency_BadElement(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	freq, err := sdr.GetTunableElementFrequency(&stub, testLogger, "IF")
	assert.NotNil(t, err)
	assert.Equal(t, "Invalid tunable element name: IF", err.Error())
	assert.Equal(t, 0.0, freq)
}

func TestSetTunableElementFrequency(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	newFreq := 50000000.
	err := sdr.SetTunableElementFrequency(&stub, testLogger, "RF", newFreq)
	assert.Nil(t, err)
	freq, err := sdr.GetTunableElementFrequency(&stub, testLogger, "RF")
	assert.Nil(t, err)
	assert.Equal(t, newFreq, freq)
}

func TestSetTunableElementFrequency_BadElement(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	newFreq := 50000000.
	err := sdr.SetTunableElementFrequency(&stub, testLogger, "IF", newFreq)
	assert.NotNil(t, err)
	assert.Equal(t, "Cannot set frequency. Invalid tunable element name: IF", err.Error())
}

func TestSetTunableElementFrequency_BadFreq(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	newFreq := -100.
	err := sdr.SetTunableElementFrequency(&stub, testLogger, "RF", newFreq)
	assert.NotNil(t, err)
	assert.Equal(t, "Cannot set frequency. Requested frequency not within element's frequency ranges", err.Error())
}
