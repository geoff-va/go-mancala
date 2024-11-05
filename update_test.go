package main

import (
	"fmt"
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

func TestHandleSwitchPlayer(t *testing.T) {
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

// Game over - P1 side is empty
func TestHandleDoneMoving_CollectRemainderP1(t *testing.T) {
	state := NewModelWithState([14]uint8{
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1,
	})
	nextState := state.HandleDoneMoving()

	assert.Equal(t, CollectRemainder, nextState, "next state")
}

// Game over - P2 side is empty
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

func TestHandleDoneMoving_StealingP1FromP2(t *testing.T) {
	state := NewModelWithState([14]uint8{
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
	})
	state.lastPlacedPit = 1
	nextState := state.HandleDoneMoving()

	assert.Equal(t, Stealing, nextState, "next state")
}

func TestHandleDoneMoving_StealingP2FromP1(t *testing.T) {
	state := NewModelWithState([14]uint8{
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
	})
	state.lastPlacedPit = 9
	state.currentPlayer = P2
	nextState := state.HandleDoneMoving()

	assert.Equal(t, Stealing, nextState, "next state")
}

func TestHandleCollectRemainder_P1Empty(t *testing.T) {
	state := NewModelWithState([14]uint8{
		5, 0, 0, 0, 0, 0, 0,
		10, 1, 1, 1, 1, 1, 1,
	})
	state.lastPlacedPit = 9
	state.currentPlayer = P1
	nextState := state.HandleCollectRemainder()

	assert.Equal(t, GameOver, nextState, "next state")
	assert.Equal(t, uint8(16), state.board.GetNumInStore(P2), "Num in P2 Store")
	assert.Equal(t, uint8(5), state.board.GetNumInStore(P1), "Num in P1 Store")
	for i := range 6 {
		assert.Equal(t, uint8(0), state.board.Get(uint8(i+1)), fmt.Sprintf("P1 side: %d", i))
		assert.Equal(t, uint8(0), state.board.Get(uint8(i+8)), fmt.Sprintf("P2 side: %d", i))
	}
}

func TestHandleCollectRemainder_P2Empty(t *testing.T) {
	state := NewModelWithState([14]uint8{
		10, 1, 1, 1, 1, 1, 1,
		5, 0, 0, 0, 0, 0, 0,
	})
	state.lastPlacedPit = 5
	state.currentPlayer = P2
	nextState := state.HandleCollectRemainder()

	assert.Equal(t, GameOver, nextState, "next state")
	assert.Equal(t, uint8(16), state.board.GetNumInStore(P1), "Num in P1 Store")
	assert.Equal(t, uint8(5), state.board.GetNumInStore(P2), "Num in P2 Store")
	for i := range 6 {
		assert.Equal(t, uint8(0), state.board.Get(uint8(i+1)), fmt.Sprintf("P1 side: %d", i))
		assert.Equal(t, uint8(0), state.board.Get(uint8(i+8)), fmt.Sprintf("P2 side: %d", i))
	}
}

func TestHandleSteal_P1FromP2(t *testing.T) {
	state := NewModelWithState([14]uint8{
		10, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 5,
	})
	state.lastPlacedPit = 1
	state.currentPlayer = P1
	nextState := state.HandleSteal()

	assert.Equal(t, SwitchPlayer, nextState, "next state")
	assert.Equal(t, uint8(0), state.board.Get(1), "Num in pit 1")
	assert.Equal(t, uint8(0), state.board.Get(13), "Num in pit 13")
	assert.Equal(t, uint8(16), state.board.GetNumInStore(P1), "Num in P1 Store")
}

func TestHandleSteal_P2FromP1(t *testing.T) {
	state := NewModelWithState([14]uint8{
		1, 1, 1, 1, 1, 1, 5,
		10, 1, 1, 1, 1, 1, 1,
	})
	state.lastPlacedPit = 8
	state.currentPlayer = P2
	nextState := state.HandleSteal()

	assert.Equal(t, SwitchPlayer, nextState, "next state")
	assert.Equal(t, uint8(0), state.board.Get(8), "Num in pit 8")
	assert.Equal(t, uint8(0), state.board.Get(6), "Num in pit 6")
	assert.Equal(t, uint8(16), state.board.GetNumInStore(P2), "Num in P2 Store")
}
