package sdr_test

import (
	"strings"
	"testing"
	"time"

	"github.com/jimorc/jsdr/internal/sdr"
	"github.com/stretchr/testify/assert"

	"github.com/jimorc/jsdr/internal/logger"
)

func TestSetupCS8Stream(t *testing.T) {
	var log strings.Builder
	testLogger := logger.New(&log)
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	assert.NotNil(t, stream)
}

func TestSetupCS8Stream_Error(t *testing.T) {
	var log strings.Builder
	testLogger := logger.New(&log)
	stub := sdr.StubDevice{Args: map[string]string{"serial": "1"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.NotNil(t, err)
	assert.Equal(t, "bad args passed to SetupCS8Stream", err.Error())
	assert.Nil(t, stream)
}

func TestCS8StreamClose(t *testing.T) {
	var log strings.Builder
	testLogger := logger.New(&log)
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	assert.NotNil(t, stream)
	err = stream.Close(testLogger)
	assert.Nil(t, err)
}

func TestGetMTU(t *testing.T) {
	var log strings.Builder
	testLogger := logger.New(&log)
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	assert.NotNil(t, stream)
	defer stream.Close(testLogger)
	mtu := stream.GetMTU(testLogger)
	assert.Equal(t, uint(10000), mtu)
}

func TestActivateCS8Stream(t *testing.T) {
	var log strings.Builder
	testLogger := logger.New(&log)
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	err = stream.Activate(testLogger, 0, 0, 0)
	assert.Nil(t, err)
	defer stream.Deactivate(testLogger, 0, 0)
}

func TestDectivateCS8Stream(t *testing.T) {
	var log strings.Builder
	testLogger := logger.New(&log)
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
	var log strings.Builder
	testLogger := logger.New(&log)
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	err = stream.Deactivate(testLogger, 0, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "attempting to deactivate a stream that is not active", err.Error())
}

func TestDeactivateCS8Stream_Error(t *testing.T) {
	var log strings.Builder
	testLogger := logger.New(&log)
	stub := sdr.StubDevice{Args: map[string]string{"serial": "2"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	err = stream.Activate(testLogger, 0, 0, 0)
	assert.Nil(t, err)
	err = stream.Deactivate(testLogger, 0, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "bad device", err.Error())
}

func TestReadCS8Stream(t *testing.T) {
	var log strings.Builder
	testLogger := logger.New(&log)
	testLogger.SetMaxLevel(logger.Debug)
	stub := sdr.StubDevice{Args: map[string]string{"serial": "3"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	err = stream.Activate(testLogger, 0, 0, 0)
	assert.Nil(t, err)
	defer stream.Deactivate(testLogger, 0, 0)
	mtu := stream.GetMTU(testLogger)
	buffer := make([][]int, 1)
	buffer[0] = make([]int, 2*mtu)
	var outputFlags [1]int
	outputFlags[0] = 0
	timeNs, numElemsRead, _ := stream.ReadCS8FromStream(testLogger, buffer, mtu, &outputFlags, 0)
	assert.True(t, timeNs > 0)
	assert.Equal(t, mtu, numElemsRead)
	assert.Equal(t, -2, buffer[0][0])
	assert.Equal(t, 0, buffer[0][1])
	assert.Equal(t, -1, buffer[0][19998])
	assert.Equal(t, -2, buffer[0][19999])
	// Wait for log to update.
	time.Sleep(100 * time.Millisecond)
	assert.True(t, strings.Contains(log.String(), "Flags ="))
	assert.True(t, strings.Contains(log.String(), "Elements Read: 10000"))
}

func TestReadCS8Stream_PartialReads(t *testing.T) {
	var log strings.Builder
	testLogger := logger.New(&log)
	testLogger.SetMaxLevel(logger.Debug)
	stub := sdr.StubDevice{Args: map[string]string{"serial": "4"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	err = stream.Activate(testLogger, 0, 0, 0)
	assert.Nil(t, err)
	defer stream.Deactivate(testLogger, 0, 0)
	mtu := stream.GetMTU(testLogger)
	buffer := make([][]int, 1)
	buffer[0] = make([]int, 2*mtu)
	var outputFlags [1]int
	timeNs, numElemsRead, _ := stream.ReadCS8FromStream(testLogger, buffer, mtu, &outputFlags, 0)
	assert.True(t, timeNs > 0)
	assert.Equal(t, mtu, numElemsRead)
	assert.Equal(t, -2, buffer[0][0])
	assert.Equal(t, 0, buffer[0][1])
	assert.Equal(t, -1, buffer[0][19998])
	assert.Equal(t, -2, buffer[0][19999])
	// Wait for log to update
	time.Sleep(100 * time.Millisecond)
	assert.True(t, strings.Contains(log.String(), "Elements Read: 5000"))
	assert.True(t, strings.Contains(log.String(), "Elements Read: 3000"))
	assert.True(t, strings.Contains(log.String(), "Elements Read: 2000"))
	assert.True(t, strings.Contains(log.String(), "Flags ="))
}

func TestReadCS8tream_NotActivated(t *testing.T) {
	var log strings.Builder
	testLogger := logger.New(&log)
	testLogger.SetMaxLevel(logger.Debug)
	stub := sdr.StubDevice{Args: map[string]string{"serial": "3"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	mtu := stream.GetMTU(testLogger)
	buffer := make([][]int, 1)
	buffer[0] = make([]int, 2*mtu)
	var outputFlags [1]int
	timeNs, numElemsRead, err := stream.ReadCS8FromStream(testLogger, buffer, mtu, &outputFlags, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "attempting to read from an inactive stream", err.Error())
	// Note: if err != nil, then there is no guarantee that timeNs and numElemsRead are valid.
	// The following tests are provided simply because this is what the test stream sets them to.
	assert.Equal(t, uint(0), timeNs)
	assert.Equal(t, uint(0), numElemsRead)
	// Wait for log to update
	time.Sleep(100 * time.Millisecond)

	assert.True(t, strings.Contains(log.String(), "Attempting to read from an inactive stream"))
}

func TestReadStreamAsCF64Data(t *testing.T) {
	var log strings.Builder
	testLogger := logger.New(&log)
	testLogger.SetMaxLevel(logger.Debug)
	stub := sdr.StubDevice{Args: map[string]string{"serial": "3"}}
	stream, err := sdr.SetupCS8Stream(&stub, testLogger)
	assert.Nil(t, err)
	defer stream.Close(testLogger)
	err = stream.Activate(testLogger, 0, 0, 0)
	assert.Nil(t, err)
	defer stream.Deactivate(testLogger, 0, 0)
	mtu := stream.GetMTU(testLogger)
	var outputFlags int
	cf64 := make([]float64, 2*mtu)
	timeNs, numElemsRead, err := stream.ReadStreamAsCF64Data(testLogger, cf64, mtu,
		&outputFlags, 0)
	assert.Nil(t, err)
	assert.Equal(t, mtu, numElemsRead)
	assert.True(t, timeNs > 0)
	assert.Equal(t, -2.0, cf64[0])
	assert.Equal(t, 0.0, cf64[1])
	assert.Equal(t, -1.0, cf64[2*mtu-2])
	assert.Equal(t, -2.0, cf64[2*mtu-1])
}
