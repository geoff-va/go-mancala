package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var moves = []struct {
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

	for _, move := range moves {
		t.Run(move.name, func(t *testing.T) {
			state := NewState()
			for range move.numMoves {
				state.MoveRight()
			}
			assert.Equal(state.selectedPit, move.expectedPit)
		})
	}

}
