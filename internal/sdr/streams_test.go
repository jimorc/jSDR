package sdr_test

import (
	"testing"

	"github.com/jimorc/jsdr/internal/sdr"
	"github.com/stretchr/testify/assert"

	"github.com/jimorc/jsdr/internal/logger"
)

func TestSetupCS8Stream(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	assert.NotNil(t, stream)
}

func TestSetupCS8Stream_Error(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.NotNil(t, err)
	assert.Equal(t, "Bad args passed to SetupCS8Stream", err.Error())
	assert.Nil(t, stream)
}

func TestCS8StreamClose(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	assert.NotNil(t, stream)
	err = stream.Close(testLogger)
	assert.Nil(t, err)
}

func TestGetMTU(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	assert.NotNil(t, stream)
	mtu := stream.GetMTU(testLogger)
	assert.Equal(t, 131072, mtu)
}
