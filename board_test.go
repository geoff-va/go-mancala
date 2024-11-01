package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveRight(t *testing.T) {
	var cases = []struct {
		name        string
		startPit    uint8
		expectedPit uint8
		player      Player
	}{
		// P1 is incrementing
		{"RightP1", 1, 2, P1},
		{"RightWrapP1", 6, 1, P1},
		// P2 is decrementing
		{"RightP2", 12, 11, P2},
		{"RightWrapP2", 8, 13, P2},
	}
	assert := assert.New(t)

	for _, move := range cases {
		t.Run(move.name, func(t *testing.T) {
			board := NewBoard()

			newPit := board.MoveRight(move.startPit, move.player)

			assert.Equal(move.expectedPit, newPit, move.name)
		})
	}
}

func TestMoveLeft(t *testing.T) {
	var cases = []struct {
		name        string
		startPit    uint8
		expectedPit uint8
		player      Player
	}{

		// P1 is decrementing
		{"LeftP1", 1, 6, P1},
		{"LeftWrapP1", 2, 1, P1},
		// P2 is incrementing
		{"LeftP2", 8, 9, P2},
		{"LeftWrapP2", 13, 8, P2},
	}
	assert := assert.New(t)

	for _, move := range cases {
		t.Run(move.name, func(t *testing.T) {
			board := NewBoard()

			newPit := board.MoveLeft(move.startPit, move.player)

			assert.Equal(move.expectedPit, newPit, move.name)
		})
	}
}

func TestSelectPit(t *testing.T) {
	board := NewBoard()
	wasInPit := board.Get(1)

	inHand := board.SelectPit(1)

	assert.Equal(t, wasInPit, inHand, "Num now in hand")
	assert.Equal(t, uint8(0), board.Get(1), "Num now in pit")
}

// P1 stealing from P2 adds both to P1's store
func TestSteal_P1(t *testing.T) {
	board := NewBoardWithOverrideState(map[uint8]uint8{1: 1})

	board.Steal(P1, 1)

	assert.Equal(t, uint8(0), board.Get(1), "Num in lastPlacedPit")
	assert.Equal(t, uint8(0), board.Get(13), "Num now in oppositePit")
	assert.Equal(t, uint8(5), board.GetNumInStore(P1), "Num in players store")
	assert.Equal(t, uint8(0), board.GetNumInStore(P2), "Num in other players store")
}

// P2 stealing from P1 adds both to P2's store
func TestSteal_P2(t *testing.T) {
	board := NewBoardWithOverrideState(map[uint8]uint8{13: 1})

	board.Steal(P2, 13)

	assert.Equal(t, uint8(0), board.Get(1), "Num in lastPlacedPit")
	assert.Equal(t, uint8(0), board.Get(13), "Num now in oppositePit")
	assert.Equal(t, uint8(5), board.GetNumInStore(P2), "Num in players store")
	assert.Equal(t, uint8(0), board.GetNumInStore(P1), "Num in other players store")
}

func TestMoveFromHandToPit_MoreInHand(t *testing.T) {
	board := NewBoard()

	inHand, lastPlacedPit := board.MoveFromHandToPit(uint8(5), 2, P1)

	assert.Equal(t, uint8(4), inHand, "inHand")
	assert.Equal(t, uint8(1), lastPlacedPit, "lastPlacedPit")
}

func TestMoveFromHandToPit_NoneLeftInHand(t *testing.T) {
	board := NewBoard()

	inHand, lastPlacedPit := board.MoveFromHandToPit(uint8(0), 2, P1)

	assert.Equal(t, uint8(0), inHand, "inHand")
	assert.Equal(t, uint8(2), lastPlacedPit, "lastPlacedPit")
}

// Skips P2 Store
func TestMoveFromHandToPit_P1NextPitIsP2Store(t *testing.T) {
	board := NewBoard()

	inHand, lastPlacedPit := board.MoveFromHandToPit(uint8(5), 8, P1)

	assert.Equal(t, uint8(4), inHand, "inHand")
	assert.Equal(t, uint8(6), lastPlacedPit, "lastPlacedPit")
}

// Skips P1 Store
func TestMoveFromHandToPit_P2NextPitIsP1Store(t *testing.T) {
	board := NewBoard()

	inHand, lastPlacedPit := board.MoveFromHandToPit(uint8(5), 1, P2)

	assert.Equal(t, uint8(4), inHand, "inHand")
	assert.Equal(t, uint8(13), lastPlacedPit, "lastPlacedPit")
}
