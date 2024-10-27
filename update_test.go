package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var rightMovesP1 = []struct {
	name        string
	numMoves    uint8
	expectedPit uint8
}{
	{"Right1", 1, 1},
	{"Right2", 2, 2},
	{"Right3", 3, 3},
	{"Right4", 4, 4},
	{"Right5", 5, 5},
	{"Right6", 6, 0},
}

func TestMoveRight_P1(t *testing.T) {
	assert := assert.New(t)

	for _, move := range rightMovesP1 {
		t.Run(move.name, func(t *testing.T) {
			state := NewState()
			for range move.numMoves {
				state.MoveRight()
			}
			assert.Equal(state.selectedPit, move.expectedPit)
		})
	}
}

var rightMovesP2 = []struct {
	name        string
	numMoves    uint8
	expectedPit uint8
}{
	{"Right1", 1, 11},
	{"Right2", 2, 10},
	{"Right3", 3, 9},
	{"Right4", 4, 8},
	{"Right5", 5, 7},
	{"Right6", 6, 12},
}

func TestMoveRight_P2(t *testing.T) {
	assert := assert.New(t)

	for _, move := range rightMovesP2 {
		t.Run(move.name, func(t *testing.T) {
			state := NewState()
			state.SwitchPlayer()
			state.selectedPit = uint8(12)
			for range move.numMoves {
				state.MoveRight()
			}
			assert.Equal(int(state.selectedPit), int(move.expectedPit))
		})
	}
}

var leftMoves = []struct {
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
	{"Left2P2", 2, 8, P2},
	{"Left3P2", 3, 9, P2},
	{"Left4P2", 4, 10, P2},
	{"Left5P2", 5, 11, P2},
	{"Left6P2", 6, 12, P2},
}

func TestMoveLeft(t *testing.T) {
	assert := assert.New(t)

	for _, move := range leftMoves {
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
