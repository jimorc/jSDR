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

func TestGetTunableElements(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	tElts := sdr.GetTunableElements(&stub, testLogger)
	assert.Equal(t, 1, len(tElts))
	assert.Equal(t, "RF", tElts[0])
}
