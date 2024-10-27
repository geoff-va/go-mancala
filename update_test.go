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
