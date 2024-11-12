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
	defer stream.Close(testLogger)
	assert.Nil(t, err)
	assert.NotNil(t, stream)
	mtu := stream.GetMTU(testLogger)
	assert.Equal(t, uint(131072), mtu)
}

func TestActivateCS8Stream(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	err = stream.Activate(testLogger, 0, 0, 0)
	assert.Nil(t, err)
	defer stream.Deactivate(testLogger, 0, 0)
}

func TestDectivateCS8Stream(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "3"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	err = stream.Activate(testLogger, 0, 0, 0)
	assert.Nil(t, err)
	err = stream.Deactivate(testLogger, 0, 0)
	assert.Nil(t, err)
}

func TestDeactivateCS8Stream_NotActive(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	err = stream.Deactivate(testLogger, 0, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "Attempting to deactivate a stream that is not active", err.Error())
}

func TestDeactivateCS8Stream_Error(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	err = stream.Activate(testLogger, 0, 0, 0)
	assert.Nil(t, err)
	err = stream.Deactivate(testLogger, 0, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "Bad device", err.Error())
}

func TestReadCS8Stream(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "3"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	err = stream.Activate(testLogger, 0, 0, 0)
	assert.Nil(t, err)
	defer stream.Deactivate(testLogger, 0, 0)
	var buffer [][]int
	var outputFlags [1]int
	timeNs, numElemsRead, err := stream.ReadCS8FromStream(testLogger, buffer, outputFlags, 0)
	assert.True(t, timeNs > 0)
	assert.Equal(t, stream.GetMTU(testLogger), numElemsRead)
}

func TestReadCS8tream_NotActivated(t *testing.T) {
	testLogger, _ := logger.NewFileLogger("stdout")
	stub := sdr.StubDevice{Args: map[string]string{"serial": "3"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	var buffer [][]int
	var outputFlags [1]int
	timeNs, numElemsRead, err := stream.ReadCS8FromStream(testLogger, buffer, outputFlags, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "Attempting to read from an inactive stream", err.Error())
	// Note: if err != nil, then there is no guarantee that timeNs and numElemsRead is valid.
	// The following tests are provided simply because this is what the test stream sets them to.
	assert.Equal(t, uint(0), timeNs)
	assert.Equal(t, uint(0), numElemsRead)
}
