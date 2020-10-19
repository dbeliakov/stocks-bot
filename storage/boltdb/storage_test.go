package boltdb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testFile          = "test.db"
	testChatID  int64 = 42
	testState         = "price"
	testSymbol1       = "AAA"
	testSymbol2       = "BBB"
)

func TestStorage_State(t *testing.T) {
	s, err := NewStorage(testFile)
	require.NoError(t, err, "Failed to create storage")
	defer func() {
		_ = os.Remove(testFile)
	}()

	err = s.Init()
	require.NoError(t, err, "Failed to init storage")

	state, err := s.GetState(testChatID)
	require.NoError(t, err, "Failed to get state")
	require.Empty(t, state, "Non empty state for new chat")

	err = s.SetState(testChatID, testState)
	require.NoError(t, err, "Failed to set state")

	state, err = s.GetState(testChatID)
	require.NoError(t, err, "Failed to get state")
	require.Equal(t, testState, state, "Not previously saved value")
}

func TestStorage_Symbols(t *testing.T) {
	s, err := NewStorage(testFile)
	require.NoError(t, err, "Failed to create storage")
	defer func() {
		_ = os.Remove(testFile)
	}()

	err = s.Init()
	require.NoError(t, err, "Failed to init storage")

	symbols, err := s.Symbols(testChatID)
	require.NoError(t, err, "Failed to get symbols")
	require.Empty(t, symbols, "Non empyty symbols for new chat")

	err = s.AddSymbol(testChatID, testSymbol1)
	require.NoError(t, err)
	symbols, err = s.Symbols(testChatID)
	require.NoError(t, err, "Failed to get symbols")
	require.Len(t, symbols, 1, "More symbols than saved")
	require.Equal(t, testSymbol1, symbols[0], "Not previously saved symbol")

	// Duplicate symbol
	err = s.AddSymbol(testChatID, testSymbol1)
	require.NoError(t, err)
	symbols, err = s.Symbols(testChatID)
	require.NoError(t, err, "Failed to get symbols")
	require.Len(t, symbols, 1, "More symbols than saved")
	require.Equal(t, testSymbol1, symbols[0], "Not previously saved symbol")

	err = s.AddSymbol(testChatID, testSymbol2)
	require.NoError(t, err)
	symbols, err = s.Symbols(testChatID)
	require.NoError(t, err, "Failed to get symbols")
	require.Len(t, symbols, 2, "More symbols than saved")
	require.Equal(t, testSymbol1, symbols[0], "Not previously saved symbol")
	require.Equal(t, testSymbol2, symbols[1], "Not previously saved symbol")

	err = s.RemoveSymbol(testChatID, testSymbol1)
	require.NoError(t, err)
	symbols, err = s.Symbols(testChatID)
	require.NoError(t, err, "Failed to get symbols")
	require.Len(t, symbols, 1, "More symbols than saved")
	require.Equal(t, testSymbol2, symbols[0], "Not previously saved symbol")
}
