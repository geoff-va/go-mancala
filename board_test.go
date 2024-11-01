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
		{"Right1P1", 1, 2, P1},
		{"Right6P2", 6, 1, P1},
		// P2 is moving right which is decrementing
		{"Right6P2", 12, 11, P2},
		{"Right1P2", 8, 13, P2},
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
