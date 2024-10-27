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
		{"Right1P1", 1, 2, P1},
		{"Right2P1", 2, 3, P1},
		{"Right3P1", 3, 4, P1},
		{"Right4P1", 4, 5, P1},
		{"Right5P1", 5, 6, P1},
		{"Right6P2", 6, 1, P1},
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

		{"Left1P1", 1, 6, P1},
		{"Left2P1", 2, 5, P1},
		{"Left3P1", 3, 4, P1},
		{"Left4P1", 4, 3, P1},
		{"Left5P1", 5, 2, P1},
		{"Left6P1", 6, 1, P1},
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
			assert.Equal(int(state.selectedPit), int(move.expectedPit))
		})
	}
}

func TestSwitchPlayer(t *testing.T) {
	assert := assert.New(t)
	state := NewState()
	assert.Equal(state.currentPlayer, P1)
	state.HandleSwitchPlayer()
	assert.Equal(state.currentPlayer, P2)
	state.HandleSwitchPlayer()
	assert.Equal(state.currentPlayer, P1)
}

func TestSelectPit(t *testing.T) {
	assert := assert.New(t)
	state := NewState()

	state.HandleMoveRight()
	nextState := state.HandleSelectPit()

	assert.Equal(state.board[state.selectedPit], uint8(0), "Num now in pit")
	assert.Equal(state.inHand, uint8(4), "inHand")
	assert.Equal(state.lastPlacedPit, state.selectedPit, "lastPlacedPit")
	assert.Equal(nextState, MovingFromHandToPit, "state")
	assert.Equal(state.lastSelectedPit[state.currentPlayer], state.selectedPit, "lastSelectedPit")

}

func TestGetOppositePit(t *testing.T) {
	var cases = []struct {
		name     string
		pit      uint8
		expected uint8
	}{
		{"Opposite1", 1, 13},
		{"Opposite2", 2, 12},
		{"Opposite3", 3, 11},
		{"Opposite4", 4, 10},
		{"Opposite5", 5, 9},
		{"Opposite6", 6, 8},
		{"Opposite8", 8, 6},
		{"Opposite9", 9, 5},
		{"Opposite10", 10, 4},
		{"Opposite11", 11, 3},
		{"Opposite12", 12, 2},
		{"Opposite13", 13, 1},
	}
	assert := assert.New(t)

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			assert.Equal(GetOppositePit(tcase.pit), tcase.expected)
		})
	}

}
