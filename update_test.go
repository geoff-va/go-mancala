package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleMoveRight(t *testing.T) {
	state := NewState()
	newState := state.HandleMoveRight()
	assert.Equal(t, SelectingPit, newState, "state")
	assert.Equal(t, uint8(2), state.selectedPit, "selectedPit")
}

func TestHandleMoveLeft(t *testing.T) {
	state := NewState()
	newState := state.HandleMoveLeft()
	assert.Equal(t, SelectingPit, newState, "state")
	assert.Equal(t, uint8(6), state.selectedPit, "selectedPit")
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

func TestHandleSelectPit(t *testing.T) {
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

func TestHandleMovingFromHandToPit_EmptyHandTurnOver(t *testing.T) {
	state := NewState()

	state.state = state.HandleSelectPit()
	state.state = state.HandleMoveFromHandToPit()
	state.state = state.HandleMoveFromHandToPit()
	state.state = state.HandleMoveFromHandToPit()
	state.state = state.HandleMoveFromHandToPit()
	nextState := state.HandleMoveFromHandToPit()

	assert.Equal(t, DoneMoving, nextState, "next state")

}
