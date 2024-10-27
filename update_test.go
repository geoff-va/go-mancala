package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveRight(t *testing.T) {
	var cases = []struct {
		name        string
		numMoves    uint8
		expectedPit uint8
		player      Player
	}{
		{"Right1P1", 1, 1, P1},
		{"Right2P1", 2, 2, P1},
		{"Right3P1", 3, 3, P1},
		{"Right4P1", 4, 4, P1},
		{"Right5P1", 5, 5, P1},
		{"Right6P2", 6, 0, P1},
		{"Right1P2", 1, 11, P2},
		{"Right2P2", 2, 10, P2},
		{"Right3P2", 3, 9, P2},
		{"Right4P2", 4, 8, P2},
		{"Right5P2", 5, 7, P2},
		{"Right6P2", 6, 12, P2},
	}
	assert := assert.New(t)

	for _, move := range cases {
		t.Run(move.name, func(t *testing.T) {
			state := NewState()
			if move.player == P2 {
				state.SwitchPlayer()
			}

			for range move.numMoves {
				state.MoveRight()
			}
			assert.Equal(state.selectedPit, move.expectedPit)
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

		{"Left1P1", 1, 5, P1},
		{"Left2P1", 2, 4, P1},
		{"Left3P1", 3, 3, P1},
		{"Left4P1", 4, 2, P1},
		{"Left5P1", 5, 1, P1},
		{"Left6P1", 6, 0, P1},
		{"Left2P1", 1, 7, P2},
		{"Left2P2", 2, 8, P2},
		{"Left3P2", 3, 9, P2},
		{"Left4P2", 4, 10, P2},
		{"Left5P2", 5, 11, P2},
		{"Left6P2", 6, 12, P2},
	}
	assert := assert.New(t)

	for _, move := range cases {
		t.Run(move.name, func(t *testing.T) {
			state := NewState()
			if move.player == P2 {
				state.SwitchPlayer()
			}

			for range move.numMoves {
				state.MoveLeft()
			}
			assert.Equal(int(state.selectedPit), int(move.expectedPit))
		})
	}
}

func TestSwitchPlayer(t *testing.T) {
	assert := assert.New(t)
	state := NewState()
	assert.Equal(state.currentPlayer, P1)
	state.SwitchPlayer()
	assert.Equal(state.currentPlayer, P2)
	state.SwitchPlayer()
	assert.Equal(state.currentPlayer, P1)
}

func TestSelectPit(t *testing.T) {
	assert := assert.New(t)
	state := NewState()

	state.MoveRight()
	state.SelectPit()

	assert.Equal(state.board[state.selectedPit], uint8(0), "Num now in pit")
	assert.Equal(state.inHand, uint8(4), "inHand")
	assert.Equal(state.selectedNum, uint8(4), "selectedNum")
	assert.Equal(state.state, MovingFromHandToPit, "state")
	assert.Equal(state.lastSelectedPit[state.currentPlayer], state.selectedPit, "lastSelectedPit")

}
