package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleMoveRight(t *testing.T) {
	state := NewModel()
	newState := state.HandleMoveRight()
	assert.Equal(t, SelectingPit, newState, "state")
	assert.Equal(t, uint8(2), state.selectedPit, "selectedPit")
}

func TestHandleMoveLeft(t *testing.T) {
	state := NewModel()
	newState := state.HandleMoveLeft()
	assert.Equal(t, SelectingPit, newState, "state")
	assert.Equal(t, uint8(6), state.selectedPit, "selectedPit")
}

func TestSwitchPlayer(t *testing.T) {
	assert := assert.New(t)
	state := NewModel()
	assert.Equal(P1, state.currentPlayer)
	state.HandleSwitchPlayer()
	assert.Equal(P2, state.currentPlayer)
	state.HandleSwitchPlayer()
	assert.Equal(P1, state.currentPlayer)
}

func TestHandleSelectPit(t *testing.T) {
	assert := assert.New(t)
	state := NewModel()

	state.HandleMoveRight()
	nextState := state.HandleSelectPit()

	assert.Equal(uint8(0), state.board.Get(state.selectedPit), uint8(0), "Num now in pit")
	assert.Equal(uint8(4), state.inHand, "inHand")
	assert.Equal(state.selectedPit, state.lastPlacedPit, "lastPlacedPit")
	assert.Equal(MovingFromHandToPit, nextState, "state")
}

func TestMovingFromHandToPit_HandNotEmpty(t *testing.T) {
	state := NewModel()

	state.state = state.HandleSelectPit()
	nextState := state.HandleMoveFromHandToPit()

	assert.Equal(t, MovingFromHandToPit, nextState, "next state")
}

func TestHandleMovingFromHandToPit_EmptyHandTurnOver(t *testing.T) {
	state := NewModel()

	state.state = state.HandleSelectPit()
	state.state = state.HandleMoveFromHandToPit()
	state.state = state.HandleMoveFromHandToPit()
	state.state = state.HandleMoveFromHandToPit()
	state.state = state.HandleMoveFromHandToPit()
	nextState := state.HandleMoveFromHandToPit()

	assert.Equal(t, DoneMoving, nextState, "next state")
}

// Turn is over
func TestHandleDoneMoving_SwitchPlayer(t *testing.T) {
	state := NewModelWithState([14]uint8{
		0, 2, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1,
	})
	state.lastPlacedPit = 1
	nextState := state.HandleDoneMoving()

	assert.Equal(t, SwitchPlayer, nextState, "next state")
}

func TestHandleDoneMoving_CollectRemainderP1(t *testing.T) {
	state := NewModelWithState([14]uint8{
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1,
	})
	nextState := state.HandleDoneMoving()

	assert.Equal(t, CollectRemainder, nextState, "next state")
}

func TestHandleDoneMoving_CollectRemainderP2(t *testing.T) {
	state := NewModelWithState([14]uint8{
		0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0,
	})
	nextState := state.HandleDoneMoving()

	assert.Equal(t, CollectRemainder, nextState, "next state")
}

// Get another turn if you land in your pit
func TestHandleDoneMoving_SelectingPitP1(t *testing.T) {
	state := NewModelWithState([14]uint8{
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
	})
	state.lastPlacedPit = 0
	nextState := state.HandleDoneMoving()

	assert.Equal(t, SelectingPit, nextState, "next state")
	assert.Equal(t, uint8(1), state.selectedPit, "selectedPit")
}

func TestHandleDoneMoving_Stealing(t *testing.T) {
}
