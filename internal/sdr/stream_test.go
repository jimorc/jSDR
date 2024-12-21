package sdr_test

import (
	"slices"
	"testing"

	"github.com/jimorc/jsdr/internal/sdr"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetStreamFormats(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	formats, err := sdr.GetStreamFormats(&stub, testLogger)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(formats))
	assert.True(t, slices.Contains(formats, "CS8"))
	assert.True(t, slices.Contains(formats, "CS16"))
	assert.True(t, slices.Contains(formats, "CF32"))
}

func TestGetStreamFormats_NoFormats(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	formats, err := sdr.GetStreamFormats(&stub, testLogger)
	assert.NotNil(t, err)
	assert.Equal(t, "no stream formats retrieved for channel 0", err.Error())
	assert.Equal(t, 0, len(formats))
}

func TestGetNativeStreamFormat(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	format, fullScale := sdr.GetNativeStreamFormat(&stub, testLogger)
	assert.Equal(t, "CS8", format)
	assert.Equal(t, 0.0, fullScale)
}
