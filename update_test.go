package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing Moving right and wrapping around
func TestHandleMoveRight(t *testing.T) {
	var cases = []struct {
		name        string
		numMoves    uint8
		expectedPit uint8
		player      Player
	}{
		{"Right1P1", 1, 2, P1},
		{"Right2P1", 2, 3, P1},
		{"Right3P1", 3, 4, P1},
		{"Right4P1", 4, 5, P1},
		{"Right5P1", 5, 6, P1},
		{"Right6P2", 6, 1, P1},
		// P2 is moving right which is decrementing
		{"Right1P2", 1, 12, P2},
		{"Right2P2", 2, 11, P2},
		{"Right3P2", 3, 10, P2},
		{"Right4P2", 4, 9, P2},
		{"Right5P2", 5, 8, P2},
		{"Right6P2", 6, 13, P2},
	}
	assert := assert.New(t)

	for _, move := range cases {
		t.Run(move.name, func(t *testing.T) {
			state := NewState()
			if move.player == P2 {
				state.HandleSwitchPlayer()
			}

			for range move.numMoves {
				state.HandleMoveRight()
			}
			assert.Equal(move.expectedPit, state.selectedPit)
		})
	}
}

func TestMoveLeft(t *testing.T) {
	var cases = []struct {
		name        string
		numMoves    uint8
		expectedPit uint8
		player      Player
	}{

		{"Left1P1", 1, 6, P1},
		{"Left2P1", 2, 5, P1},
		{"Left3P1", 3, 4, P1},
		{"Left4P1", 4, 3, P1},
		{"Left5P1", 5, 2, P1},
		{"Left6P1", 6, 1, P1},
		// P2 is moving right which is incrementing
		{"Left2P1", 1, 8, P2},
		{"Left2P2", 2, 9, P2},
		{"Left3P2", 3, 10, P2},
		{"Left4P2", 4, 11, P2},
		{"Left5P2", 5, 12, P2},
		{"Left6P2", 6, 13, P2},
	}
	assert := assert.New(t)

	for _, move := range cases {
		t.Run(move.name, func(t *testing.T) {
			state := NewState()
			if move.player == P2 {
				state.HandleSwitchPlayer()
			}

			for range move.numMoves {
				state.HandleMoveLeft()
			}
			assert.Equal(move.expectedPit, state.selectedPit)
		})
	}
}

func TestSwitchPlayer(t *testing.T) {
	assert := assert.New(t)
	state := NewState()
	assert.Equal(P1, state.currentPlayer)
	state.HandleSwitchPlayer()
	assert.Equal(P2, state.currentPlayer)
	state.HandleSwitchPlayer()
	assert.Equal(P1, state.currentPlayer)
}

func TestSelectPit(t *testing.T) {
	assert := assert.New(t)
	state := NewState()

	state.HandleMoveRight()
	nextState := state.HandleSelectPit()

	assert.Equal(uint8(0), state.board.Get(state.selectedPit), uint8(0), "Num now in pit")
	assert.Equal(uint8(4), state.inHand, "inHand")
	assert.Equal(state.selectedPit, state.lastPlacedPit, "lastPlacedPit")
	assert.Equal(MovingFromHandToPit, nextState, "state")

}
func TestMovingFromHandToPit_HandNotEmpty(t *testing.T) {
	state := NewState()

	state.state = state.HandleSelectPit()
	nextState := state.HandleMoveFromHandToPit()

	assert.Equal(t, MovingFromHandToPit, nextState, "next state")

}

func TestMovingFromHandToPit_EmptyHandTurnOver(t *testing.T) {
	state := NewState()

	state.state = state.HandleSelectPit()
	state.state = state.HandleMoveFromHandToPit()
	state.state = state.HandleMoveFromHandToPit()
	state.state = state.HandleMoveFromHandToPit()
	state.state = state.HandleMoveFromHandToPit()
	nextState := state.HandleMoveFromHandToPit()

	assert.Equal(t, SwitchPlayer, nextState, "next state")

}
